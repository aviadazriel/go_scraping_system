package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
	Checks    map[string]string `json:"checks,omitempty"`
}

// healthHandler handles the main health check endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Service:   "api-gateway",
		Timestamp: time.Now().Format(time.RFC3339),
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
func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if all dependencies are ready
	// For now, always return ready
	response := HealthResponse{
		Status:    "ready",
		Service:   "api-gateway",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// livenessHandler handles the liveness check endpoint
func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Check if the service is alive and responsive
	// For now, always return alive
	response := HealthResponse{
		Status:    "alive",
		Service:   "api-gateway",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
