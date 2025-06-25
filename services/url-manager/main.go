package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go_scraping_project/services/url-manager/repositories"
	"go_scraping_project/services/url-manager/services"
	"go_scraping_project/shared/config"
	"go_scraping_project/shared/database"
	"go_scraping_project/shared/kafka"

	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration using shared config loader
	loader := config.NewLoader()
	if err := loader.LoadServiceConfig("url-manager"); err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Load environment variables
	loader.LoadFromEnv()

	// Initialize logger based on config
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level from config
	logLevel := loader.GetString("logging.level")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set DATABASE_URL environment variable for shared database package
	dbHost := loader.GetString("database.host")
	dbPort := loader.GetInt("database.port")
	dbUser := loader.GetString("database.user")
	dbPassword := loader.GetString("database.password")
	dbName := loader.GetString("database.db_name")
	dbSSLMode := loader.GetString("database.ssl_mode")

	if dbName == "" {
		dbName = "scraping_db" // fallback
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable" // fallback
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	os.Setenv("DATABASE_URL", databaseURL)

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Initialize sqlc-generated database queries
	queries := database.New(db)

	// Initialize Kafka producer using config
	kafkaBrokers := loader.GetStringSlice("kafka.brokers")
	if len(kafkaBrokers) == 0 {
		kafkaBrokers = []string{"localhost:9092"} // fallback
	}

	producer, err := kafka.NewProducer(kafkaBrokers, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Kafka producer")
	}
	defer producer.Close()

	// Initialize URL repository
	urlRepo := repositories.NewURLRepository(queries, logger)

	// Initialize URL scheduler service
	scheduler := services.NewURLSchedulerService(urlRepo, producer, logger)

	// Start scheduler
	logger.Info("Starting URL scheduler service")
	if err := scheduler.Start(context.Background()); err != nil {
		logger.WithError(err).Fatal("Failed to start scheduler")
	}

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down URL Manager...")

	// Stop scheduler
	scheduler.Stop()

	logger.Info("URL Manager exited")
}
