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

// getDatabaseURL returns the database connection URL, prioritizing environment variables over config
func getDatabaseURL(cfg *config.Loader) string {
	// Check if DATABASE_URL is set in environment
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	// Build from config if not in environment
	dbHost := cfg.GetString("database.host")
	dbPort := cfg.GetInt("database.port")
	dbUser := cfg.GetString("database.user")
	dbPassword := cfg.GetString("database.password")
	dbName := cfg.GetString("database.database")
	dbSSLMode := cfg.GetString("database.ssl_mode")

	// Set defaults
	if dbName == "" {
		dbName = "scraping_db"
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
}

// getLogger initializes and configures a logger based on configuration
func getLogger(cfg *config.Loader) *logrus.Logger {
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

	return logger
}

// getServerConfig returns server configuration with fallbacks
func getServerConfig(cfg *config.Loader) (int, time.Duration, time.Duration, time.Duration) {
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

	return port, readTimeoutDuration, writeTimeoutDuration, idleTimeoutDuration
}

// createServer creates and configures the HTTP server
func createServer(handler http.Handler, port int, readTimeout, writeTimeout, idleTimeout time.Duration) *http.Server {
	return &http.Server{
		Addr:         ":" + strconv.Itoa(port),
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}
}

// startServer starts the HTTP server in a goroutine
func startServer(server *http.Server, logger *logrus.Logger, port int) {
	go func() {
		logger.Infof("Starting API Gateway server on :%d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()
}

// waitForShutdown waits for interrupt signal and gracefully shuts down the server
func waitForShutdown(server *http.Server, logger *logrus.Logger) {
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
	logger := getLogger(cfg)

	// Get database URL and set environment variable
	databaseURL := getDatabaseURL(cfg)
	os.Setenv("DATABASE_URL", databaseURL)

	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer db.Close()

	// Initialize sqlc-generated database queries
	queries := database.New(db)

	// Initialize router
	router := handlers.NewRouter(logger, queries)
	handler := handlers.SetupRoutes(router)

	// Get server configuration
	port, readTimeout, writeTimeout, idleTimeout := getServerConfig(cfg)

	// Create HTTP server
	server := createServer(handler, port, readTimeout, writeTimeout, idleTimeout)

	// Start server
	startServer(server, logger, port)

	// Wait for shutdown
	waitForShutdown(server, logger)
}
