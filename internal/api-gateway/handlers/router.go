package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router handles all route registration for the API Gateway
type Router struct {
	router *mux.Router
	logger *logrus.Logger

	// Handlers
	urlHandler     *URLHandler
	dataHandler    *DataHandler
	metricsHandler *MetricsHandler
	adminHandler   *AdminHandler
}

// NewRouter creates a new router with all handlers
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
func (r *Router) setupDataRoutes(apiV1 *mux.Router) {
	dataRoutes := apiV1.PathPrefix("/data").Subrouter()

	dataRoutes.HandleFunc("", r.dataHandler.ListData).Methods("GET")
	dataRoutes.HandleFunc("/{url_id}", r.dataHandler.GetDataByURL).Methods("GET")
	dataRoutes.HandleFunc("/export", r.dataHandler.ExportData).Methods("GET")
}

// setupMetricsRoutes configures metrics routes
func (r *Router) setupMetricsRoutes(apiV1 *mux.Router) {
	metricsRoutes := apiV1.PathPrefix("/metrics").Subrouter()

	metricsRoutes.HandleFunc("/urls/{id}", r.metricsHandler.GetURLMetrics).Methods("GET")
	metricsRoutes.HandleFunc("/system", r.metricsHandler.GetSystemMetrics).Methods("GET")
}

// setupAdminRoutes configures admin routes
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
func (r *Router) GetRouter() *mux.Router {
	return r.router
}
