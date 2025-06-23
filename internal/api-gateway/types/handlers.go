package types

import (
	"go_scraping_project/internal/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router handles all route registration for the API Gateway
// It provides a centralized way to organize and configure all HTTP routes
// for the web scraping system, including middleware setup and route grouping.
type Router struct {
	Router *mux.Router
	Logger *logrus.Logger
	DB     *database.Queries

	// Handlers
	URLHandler     *URLHandler     // Handles URL management endpoints
	DataHandler    *DataHandler    // Handles data retrieval endpoints
	MetricsHandler *MetricsHandler // Handles metrics and monitoring endpoints
	AdminHandler   *AdminHandler   // Handles admin and system management endpoints
}
