package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/pkg/database"
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
	}).Info("Starting URL Manager Service")

	// Initialize database
	db, err := database.NewPostgresDB(cfg, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Errorf("Failed to close database connection: %v", err)
		}
	}()

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

	// Initialize Kafka consumer
	consumer, err := kafka.NewConsumer(cfg, log)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Errorf("Failed to close Kafka consumer: %v", err)
		}
	}()

	// Initialize HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      createHandler(cfg, db, producer, log),
	}

	// Start HTTP server in a goroutine
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start Kafka consumer in a goroutine
	go func() {
		topics := []string{"retry", "dead-letter"}
		log.WithField("topics", topics).Info("Starting Kafka consumer")
		if err := consumer.Consume(topics); err != nil {
			log.Errorf("Kafka consumer error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down URL Manager Service...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("URL Manager Service exited")
}

// createHandler creates the HTTP handler for the URL Manager service
func createHandler(cfg *config.Config, db *database.PostgresDB, producer *kafka.Producer, log *logrus.Logger) http.Handler {
	// TODO: Initialize repositories and services
	// TODO: Create HTTP handlers
	// TODO: Add middleware (logging, metrics, tracing, etc.)

	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"url-manager"}`))
	})

	// TODO: Add other endpoints
	// - POST /api/v1/urls - Create URL
	// - GET /api/v1/urls - List URLs
	// - GET /api/v1/urls/{id} - Get URL
	// - PUT /api/v1/urls/{id} - Update URL
	// - DELETE /api/v1/urls/{id} - Delete URL
	// - POST /api/v1/urls/{id}/scrape - Trigger immediate scrape

	return mux
}
