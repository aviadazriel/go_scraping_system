package handlers

import (
	"net/http"

	"go_scraping_project/services/api-gateway/types"
	"go_scraping_project/shared/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// NewRouter creates a new router with all handlers and database queries
// This function initializes the router and all handler instances with their
// required database queries, setting up the complete routing structure
// for the API Gateway.
//
// Parameters:
//   - logger: Structured logger for request logging and error handling
//   - db: sqlc-generated database queries for data persistence
//
// Returns:
//   - *types.Router: Configured router instance ready for route setup
func NewRouter(logger *logrus.Logger, db *database.Queries) *types.Router {
	router := mux.NewRouter()

	// Initialize handlers with database queries
	urlHandler := types.NewURLHandler(logger, db)
	dataHandler := types.NewDataHandler(logger)
	metricsHandler := types.NewMetricsHandler(logger)
	adminHandler := types.NewAdminHandler(logger)

	return &types.Router{
		Router:         router,
		Logger:         logger,
		DB:             db,
		URLHandler:     urlHandler,
		DataHandler:    dataHandler,
		MetricsHandler: metricsHandler,
		AdminHandler:   adminHandler,
	}
}

// SetupRoutes configures all routes and middleware
//
// Purpose: Sets up the complete routing structure for the API Gateway,
// including middleware configuration, health check endpoints, and API routes.
// This method organizes routes by functionality and applies appropriate
// middleware for logging, CORS, and error handling.
//
// Returns:
//   - http.Handler: Configured HTTP handler ready for server use
//
// Route Structure:
//   - Health endpoints: /health, /ready, /live
//   - API v1 endpoints: /api/v1/*
//   - URL management: /api/v1/urls/*
//   - Data retrieval: /api/v1/data/*
//   - Metrics: /api/v1/metrics/*
//   - Admin: /api/v1/admin/*
//
// Middleware Applied:
//   - Logging middleware for request tracking
//   - CORS middleware for cross-origin support
//   - Recovery middleware for panic handling
func SetupRoutes(router *types.Router) http.Handler {
	// Add middleware
	router.Router.Use(loggingMiddleware(router.Logger))
	router.Router.Use(corsMiddleware())
	router.Router.Use(recoveryMiddleware(router.Logger))

	// Health check endpoints
	router.Router.HandleFunc("/health", healthHandler).Methods("GET")
	router.Router.HandleFunc("/ready", readinessHandler).Methods("GET")
	router.Router.HandleFunc("/live", livenessHandler).Methods("GET")

	// API v1 routes
	apiV1 := router.Router.PathPrefix("/api/v1").Subrouter()

	// Setup route groups
	setupURLRoutes(apiV1, router.URLHandler)
	setupDataRoutes(apiV1, router.DataHandler)
	setupMetricsRoutes(apiV1, router.MetricsHandler)
	setupAdminRoutes(apiV1, router.AdminHandler)

	return router.Router
}

// setupURLRoutes configures URL management routes
//
// Purpose: Sets up all routes related to URL management, including
// CRUD operations, status monitoring, and manual scraping triggers.
//
// Routes Configured:
//   - POST /api/v1/urls - Create a new URL
//   - GET /api/v1/urls - List all URLs (with pagination)
//   - GET /api/v1/urls/{id} - Get specific URL details
//   - PUT /api/v1/urls/{id} - Update URL configuration
//   - DELETE /api/v1/urls/{id} - Delete a URL
//   - POST /api/v1/urls/{id}/scrape - Trigger manual scraping
//   - GET /api/v1/urls/{id}/status - Get URL status information
//
// Parameters:
//   - apiV1: Subrouter for API v1 endpoints
//   - urlHandler: URL handler instance
func setupURLRoutes(apiV1 *mux.Router, urlHandler *types.URLHandler) {
	urlRoutes := apiV1.PathPrefix("/urls").Subrouter()

	urlRoutes.HandleFunc("", urlHandler.CreateURL).Methods("POST")
	urlRoutes.HandleFunc("", urlHandler.ListURLs).Methods("GET")
	urlRoutes.HandleFunc("/{id}", urlHandler.GetURL).Methods("GET")
	urlRoutes.HandleFunc("/{id}", urlHandler.UpdateURL).Methods("PUT")
	urlRoutes.HandleFunc("/{id}", urlHandler.DeleteURL).Methods("DELETE")
	urlRoutes.HandleFunc("/{id}/scrape", urlHandler.TriggerScrape).Methods("POST")
	urlRoutes.HandleFunc("/{id}/status", urlHandler.GetURLStatus).Methods("GET")
}

// setupDataRoutes configures data retrieval routes
//
// Purpose: Sets up all routes related to data retrieval and export,
// providing access to scraped and parsed data with filtering options.
//
// Routes Configured:
//   - GET /api/v1/data - List scraped data (with filtering and pagination)
//   - GET /api/v1/data/{url_id} - Get data for specific URL
//   - GET /api/v1/data/export - Export data in various formats
//
// Parameters:
//   - apiV1: Subrouter for API v1 endpoints
//   - dataHandler: Data handler instance
func setupDataRoutes(apiV1 *mux.Router, dataHandler *types.DataHandler) {
	dataRoutes := apiV1.PathPrefix("/data").Subrouter()

	dataRoutes.HandleFunc("", dataHandler.ListData).Methods("GET")
	dataRoutes.HandleFunc("/{url_id}", dataHandler.GetDataByURL).Methods("GET")
	dataRoutes.HandleFunc("/export", dataHandler.ExportData).Methods("GET")
}

// setupMetricsRoutes configures metrics routes
//
// Purpose: Sets up all routes related to system metrics and monitoring,
// providing insights into system performance and URL-specific statistics.
//
// Routes Configured:
//   - GET /api/v1/metrics/urls/{id} - Get metrics for specific URL
//   - GET /api/v1/metrics/system - Get system-wide metrics
//
// Parameters:
//   - apiV1: Subrouter for API v1 endpoints
//   - metricsHandler: Metrics handler instance
func setupMetricsRoutes(apiV1 *mux.Router, metricsHandler *types.MetricsHandler) {
	metricsRoutes := apiV1.PathPrefix("/metrics").Subrouter()

	metricsRoutes.HandleFunc("/urls/{id}", metricsHandler.GetURLMetrics).Methods("GET")
	metricsRoutes.HandleFunc("/system", metricsHandler.GetSystemMetrics).Methods("GET")
}

// setupAdminRoutes configures admin routes
//
// Purpose: Sets up all routes related to system administration,
// including dead letter queue management and comprehensive health monitoring.
//
// Routes Configured:
//   - GET /api/v1/admin/dead-letter - List dead letter messages
//   - POST /api/v1/admin/dead-letter/bulk-retry - Bulk retry failed messages
//   - POST /api/v1/admin/dead-letter/{id}/retry - Retry specific message
//   - DELETE /api/v1/admin/dead-letter/{id} - Delete dead letter message
//   - GET /api/v1/admin/health - Get comprehensive system health
//
// Parameters:
//   - apiV1: Subrouter for API v1 endpoints
//   - adminHandler: Admin handler instance
func setupAdminRoutes(apiV1 *mux.Router, adminHandler *types.AdminHandler) {
	adminRoutes := apiV1.PathPrefix("/admin").Subrouter()

	// Dead letter queue management
	adminRoutes.HandleFunc("/dead-letter", adminHandler.ListDeadLetterMessages).Methods("GET")
	adminRoutes.HandleFunc("/dead-letter/bulk-retry", adminHandler.BulkRetryDeadLetterMessages).Methods("POST")
	adminRoutes.HandleFunc("/dead-letter/{id}/retry", adminHandler.RetryDeadLetterMessage).Methods("POST")
	adminRoutes.HandleFunc("/dead-letter/{id}", adminHandler.DeleteDeadLetterMessage).Methods("DELETE")

	// System health
	adminRoutes.HandleFunc("/health", adminHandler.GetSystemHealth).Methods("GET")
}
