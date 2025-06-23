package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router handles all route registration for the API Gateway
// It provides a centralized way to organize and configure all HTTP routes
// for the web scraping system, including middleware setup and route grouping.
type Router struct {
	router *mux.Router
	logger *logrus.Logger

	// Handlers
	urlHandler     *URLHandler     // Handles URL management endpoints
	dataHandler    *DataHandler    // Handles data retrieval endpoints
	metricsHandler *MetricsHandler // Handles metrics and monitoring endpoints
	adminHandler   *AdminHandler   // Handles admin and system management endpoints
}

// NewRouter creates a new router with all handlers
// This function initializes the router and all handler instances,
// setting up the complete routing structure for the API Gateway.
//
// Parameters:
//   - logger: Structured logger for request logging and error handling
//
// Returns:
//   - *Router: Configured router instance ready for route setup
func NewRouter(logger *logrus.Logger) *Router {
	router := mux.NewRouter()

	// Initialize handlers
	urlHandler := NewURLHandler(logger)
	dataHandler := NewDataHandler(logger)
	metricsHandler := NewMetricsHandler(logger)
	adminHandler := NewAdminHandler(logger)

	return &Router{
		router:         router,
		logger:         logger,
		urlHandler:     urlHandler,
		dataHandler:    dataHandler,
		metricsHandler: metricsHandler,
		adminHandler:   adminHandler,
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
func (r *Router) SetupRoutes() http.Handler {
	// Add middleware
	r.router.Use(loggingMiddleware(r.logger))
	r.router.Use(corsMiddleware())
	r.router.Use(recoveryMiddleware(r.logger))

	// Health check endpoints
	r.router.HandleFunc("/health", healthHandler).Methods("GET")
	r.router.HandleFunc("/ready", readinessHandler).Methods("GET")
	r.router.HandleFunc("/live", livenessHandler).Methods("GET")

	// API v1 routes
	apiV1 := r.router.PathPrefix("/api/v1").Subrouter()

	// Setup route groups
	r.setupURLRoutes(apiV1)
	r.setupDataRoutes(apiV1)
	r.setupMetricsRoutes(apiV1)
	r.setupAdminRoutes(apiV1)

	return r.router
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
func (r *Router) setupURLRoutes(apiV1 *mux.Router) {
	urlRoutes := apiV1.PathPrefix("/urls").Subrouter()

	urlRoutes.HandleFunc("", r.urlHandler.CreateURL).Methods("POST")
	urlRoutes.HandleFunc("", r.urlHandler.ListURLs).Methods("GET")
	urlRoutes.HandleFunc("/{id}", r.urlHandler.GetURL).Methods("GET")
	urlRoutes.HandleFunc("/{id}", r.urlHandler.UpdateURL).Methods("PUT")
	urlRoutes.HandleFunc("/{id}", r.urlHandler.DeleteURL).Methods("DELETE")
	urlRoutes.HandleFunc("/{id}/scrape", r.urlHandler.TriggerScrape).Methods("POST")
	urlRoutes.HandleFunc("/{id}/status", r.urlHandler.GetURLStatus).Methods("GET")
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
func (r *Router) setupDataRoutes(apiV1 *mux.Router) {
	dataRoutes := apiV1.PathPrefix("/data").Subrouter()

	dataRoutes.HandleFunc("", r.dataHandler.ListData).Methods("GET")
	dataRoutes.HandleFunc("/{url_id}", r.dataHandler.GetDataByURL).Methods("GET")
	dataRoutes.HandleFunc("/export", r.dataHandler.ExportData).Methods("GET")
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
func (r *Router) setupMetricsRoutes(apiV1 *mux.Router) {
	metricsRoutes := apiV1.PathPrefix("/metrics").Subrouter()

	metricsRoutes.HandleFunc("/urls/{id}", r.metricsHandler.GetURLMetrics).Methods("GET")
	metricsRoutes.HandleFunc("/system", r.metricsHandler.GetSystemMetrics).Methods("GET")
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
func (r *Router) setupAdminRoutes(apiV1 *mux.Router) {
	adminRoutes := apiV1.PathPrefix("/admin").Subrouter()

	// Dead letter queue management
	adminRoutes.HandleFunc("/dead-letter", r.adminHandler.ListDeadLetterMessages).Methods("GET")
	adminRoutes.HandleFunc("/dead-letter/bulk-retry", r.adminHandler.BulkRetryDeadLetterMessages).Methods("POST")
	adminRoutes.HandleFunc("/dead-letter/{id}/retry", r.adminHandler.RetryDeadLetterMessage).Methods("POST")
	adminRoutes.HandleFunc("/dead-letter/{id}", r.adminHandler.DeleteDeadLetterMessage).Methods("DELETE")

	// System health
	adminRoutes.HandleFunc("/health", r.adminHandler.GetSystemHealth).Methods("GET")
}

// GetRouter returns the underlying mux.Router for additional customization if needed
//
// Purpose: Provides access to the underlying gorilla/mux router for advanced
// customization, such as adding custom middleware, route handlers, or
// implementing additional routing logic not covered by the standard setup.
//
// Returns:
//   - *mux.Router: The underlying gorilla/mux router instance
//
// Example Usage:
//
//	router := handlerRouter.GetRouter()
//	router.Use(customMiddleware)
//	router.HandleFunc("/custom", customHandler)
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
