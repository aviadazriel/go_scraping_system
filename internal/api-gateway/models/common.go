package models

import (
	"net/http"
)

// ValidationError represents a validation error with field-specific information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return e.Message
}

// responseWriter wraps http.ResponseWriter to capture status code
// This is used by middleware for logging and metrics
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code for logging purposes
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
