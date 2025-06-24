package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/internal/database"
	"go_scraping_project/internal/url-manager/repositories"
	"go_scraping_project/internal/url-manager/services"
	pkgdb "go_scraping_project/pkg/database"
	"go_scraping_project/pkg/kafka"
	"go_scraping_project/pkg/observability"

	"github.com/sirupsen/logrus"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := observability.NewLogger(cfg)
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"version":    Version,
		"build_time": BuildTime,
		"service":    "url-manager",
	}).Info("Starting URL Manager Background Service")

	// Initialize database
	db, err := pkgdb.NewPostgresDB(cfg, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("Failed to close database connection: %v", err)
		}
	}()

	// Initialize database queries using sqlc
	queries := database.New(db.GetDB())

	// Initialize URL repository
	urlRepo := repositories.NewURLRepository(queries, log)

	// Initialize Kafka producer
	producer, err := kafka.NewProducer(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Errorf("Failed to close Kafka producer: %v", err)
		}
	}()

	// Wrap Kafka producer to match domain interface
	kafkaProducer := kafka.NewKafkaProducerWrapper(producer)

	// Initialize URL scheduler service
	schedulerService := services.NewURLSchedulerService(urlRepo, kafkaProducer, log)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the scheduler service
	if err := schedulerService.Start(ctx); err != nil {
		log.Fatalf("Failed to start scheduler service: %v", err)
	}
	defer func() {
		if err := schedulerService.Stop(); err != nil {
			log.Errorf("Failed to stop scheduler service: %v", err)
		}
	}()

	log.Info("URL Manager Background Service is running. Press Ctrl+C to stop.")

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down URL Manager Background Service...")

	// Give the scheduler time to finish current operations
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// The scheduler will stop when the context is cancelled
	cancel()

	// Wait for shutdown timeout or completion
	<-shutdownCtx.Done()

	log.Info("URL Manager Background Service exited")
}
