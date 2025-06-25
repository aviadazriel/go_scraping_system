package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go_scraping_project/services/api-gateway/handlers"
	"go_scraping_project/shared/config"
	"go_scraping_project/shared/database"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables from .env if present (for backward compatibility)
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.NewLoader()
	if err := cfg.LoadServiceConfig("api-gateway"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Enable environment variable overrides
	cfg.LoadFromEnv()

	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level from config
	logLevel := cfg.GetString("logging.level")
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Initialize database connection
	db, err := database.ConnectWithConfig(cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Initialize sqlc-generated database queries
	queries := database.New(db)

	// Initialize router
	router := handlers.NewRouter(logger, queries)
	handler := handlers.SetupRoutes(router)

	// Get server configuration from config
	port := cfg.GetInt("server.port")
	if port == 0 {
		port = 8080 // fallback
	}

	readTimeout := cfg.GetString("server.read_timeout")
	writeTimeout := cfg.GetString("server.write_timeout")
	idleTimeout := cfg.GetString("server.idle_timeout")

	// Parse timeouts (with fallbacks)
	readTimeoutDuration, _ := time.ParseDuration(readTimeout)
	if readTimeoutDuration == 0 {
		readTimeoutDuration = 30 * time.Second
	}

	writeTimeoutDuration, _ := time.ParseDuration(writeTimeout)
	if writeTimeoutDuration == 0 {
		writeTimeoutDuration = 30 * time.Second
	}

	idleTimeoutDuration, _ := time.ParseDuration(idleTimeout)
	if idleTimeoutDuration == 0 {
		idleTimeoutDuration = 60 * time.Second
	}

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      handler,
		ReadTimeout:  readTimeoutDuration,
		WriteTimeout: writeTimeoutDuration,
		IdleTimeout:  idleTimeoutDuration,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting API Gateway server on :%d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
}
