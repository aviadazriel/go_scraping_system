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
	"go_scraping_project/internal/handlers"
	"go_scraping_project/pkg/observability"

	"github.com/gorilla/mux"
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
		"service":    "api-gateway",
	}).Info("Starting API Gateway Service")

	// Create router
	router := createRouter(cfg, log)

	// Initialize HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		Handler:      router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.WithField("port", cfg.Server.Port).Info("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Gateway Service...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	}

	log.Info("API Gateway Service exited")
}

// createRouter creates the HTTP router with all routes and middleware
func createRouter(cfg *config.Config, log *logrus.Logger) http.Handler {
	router := mux.NewRouter()

	// Add middleware
	router.Use(loggingMiddleware(log))
	router.Use(corsMiddleware())
	router.Use(recoveryMiddleware(log))

	// Initialize handlers
	urlHandler := handlers.NewURLHandler(log)
	dataHandler := handlers.NewDataHandler(log)

	// Health check endpoint
	router.HandleFunc("/health", healthHandler).Methods("GET")

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// URL management routes
	urlRoutes := apiV1.PathPrefix("/urls").Subrouter()
	urlRoutes.HandleFunc("", urlHandler.CreateURL).Methods("POST")
	urlRoutes.HandleFunc("", urlHandler.ListURLs).Methods("GET")
	urlRoutes.HandleFunc("/{id}", urlHandler.GetURL).Methods("GET")
	urlRoutes.HandleFunc("/{id}", urlHandler.UpdateURL).Methods("PUT")
	urlRoutes.HandleFunc("/{id}", urlHandler.DeleteURL).Methods("DELETE")
	urlRoutes.HandleFunc("/{id}/scrape", urlHandler.TriggerScrape).Methods("POST")
	urlRoutes.HandleFunc("/{id}/status", urlHandler.GetURLStatus).Methods("GET")

	// Data retrieval routes
	dataRoutes := apiV1.PathPrefix("/data").Subrouter()
	dataRoutes.HandleFunc("", dataHandler.ListData).Methods("GET")
	dataRoutes.HandleFunc("/{url_id}", dataHandler.GetDataByURL).Methods("GET")
	dataRoutes.HandleFunc("/export", dataHandler.ExportData).Methods("GET")

	// Metrics routes
	metricsRoutes := apiV1.PathPrefix("/metrics").Subrouter()
	metricsRoutes.HandleFunc("/urls/{id}", getURLMetricsHandler).Methods("GET")
	metricsRoutes.HandleFunc("/system", getSystemMetricsHandler).Methods("GET")

	// Admin routes
	adminRoutes := apiV1.PathPrefix("/admin").Subrouter()
	adminRoutes.HandleFunc("/dead-letter", listDeadLetterMessagesHandler).Methods("GET")
	adminRoutes.HandleFunc("/dead-letter/{id}/retry", retryDeadLetterMessageHandler).Methods("POST")
	adminRoutes.HandleFunc("/dead-letter/{id}", deleteDeadLetterMessageHandler).Methods("DELETE")

	return router
}

// Health check handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"api-gateway","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`))
}

// Metrics handlers
func getURLMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement URL metrics retrieval
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
}

func getSystemMetricsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement system metrics retrieval
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
}

// Admin handlers
func listDeadLetterMessagesHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement dead letter message listing
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
}

func retryDeadLetterMessageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement dead letter message retry
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
}

func deleteDeadLetterMessageHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement dead letter message deletion
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
}

// Middleware functions
func loggingMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			log.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     wrapped.statusCode,
				"duration":   duration,
				"user_agent": r.UserAgent(),
				"remote_ip":  r.RemoteAddr,
			}).Info("HTTP Request")
		})
	}
}

func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func recoveryMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.WithFields(logrus.Fields{
						"error":  err,
						"path":   r.URL.Path,
						"method": r.Method,
					}).Error("Panic recovered")

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(`{"error":"Internal server error"}`))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
