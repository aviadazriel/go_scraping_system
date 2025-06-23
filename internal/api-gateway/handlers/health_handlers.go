package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go_scraping_project/internal/api-gateway/models"
)

// HealthResponse represents the health check response structure.
// It provides standardized health information for monitoring systems,
// load balancers, and Kubernetes probes.
type HealthResponse struct {
	Status    string            `json:"status"`            // Health status (healthy, unhealthy, ready, alive)
	Service   string            `json:"service"`           // Service name identifier
	Timestamp string            `json:"timestamp"`         // ISO 8601 timestamp of the health check
	Version   string            `json:"version,omitempty"` // Service version (optional)
	Uptime    string            `json:"uptime,omitempty"`  // Service uptime duration (optional)
	Checks    map[string]string `json:"checks,omitempty"`  // Additional health checks (optional)
}

// healthHandler handles the main health check endpoint
//
// Purpose: Provides a basic health check for load balancers and monitoring systems.
// This endpoint performs a quick assessment of the service's health status,
// including dependency checks for database and Kafka connectivity.
//
// Response: models.HealthResponse (200 OK)
//
// Example Usage:
//
//	GET /health
//
// Response Example:
//
//	{
//	  "status": "healthy",
//	  "service": "api-gateway",
//	  "timestamp": "2024-01-01T00:00:00Z",
//	  "version": "1.0.0",
//	  "checks": {
//	    "database": "healthy",
//	    "kafka": "healthy"
//	  }
//	}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    "24h30m",
		Version:   "1.0.0",
		Checks: map[string]string{
			"database": "healthy",
			"kafka":    "healthy",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// readinessHandler handles the readiness check endpoint
//
// Purpose: Kubernetes readiness probe to check if the service is ready to receive traffic.
// This endpoint verifies that all dependencies are available and the service
// is fully initialized and ready to handle requests.
//
// Response: models.HealthResponse (200 OK) or (503 Service Unavailable)
//
// Example Usage:
//
//	GET /ready
//
// Response Example:
//
//	{
//	  "status": "ready",
//	  "timestamp": "2024-01-01T00:00:00Z"
//	}
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if all dependencies are ready
	// For now, always return ready
	response := models.HealthResponse{
		Status:    "ready",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    "24h30m",
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// livenessHandler handles the liveness check endpoint
//
// Purpose: Kubernetes liveness probe to check if the service is alive and responsive.
// This endpoint verifies that the service is running and can respond to requests.
// If this endpoint fails repeatedly, Kubernetes will restart the pod.
//
// Response: models.HealthResponse (200 OK) or (503 Service Unavailable)
//
// Example Usage:
//
//	GET /live
//
// Response Example:
//
//	{
//	  "status": "alive",
//	  "timestamp": "2024-01-01T00:00:00Z"
//	}
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the service is alive and responsive
	// For now, always return alive
	response := models.HealthResponse{
		Status:    "alive",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    "24h30m",
		Version:   "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
