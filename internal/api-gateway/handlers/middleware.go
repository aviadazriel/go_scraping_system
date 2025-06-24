package handlers

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// responseWriter wraps http.ResponseWriter to capture status code
// This wrapper allows middleware to capture the HTTP status code
// that was set by the handler for logging purposes.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code before calling the original WriteHeader
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs all HTTP requests with structured logging
//
// Purpose: Provides comprehensive request logging for all HTTP requests.
// This middleware captures request details including method, path, status code,
// response time, user agent, and remote IP address for monitoring and debugging.
//
// Features:
//   - Structured logging with consistent format
//   - Request duration measurement
//   - Status code capture
//   - User agent and IP address logging
//
// Example Usage:
//
//	router.Use(loggingMiddleware(logger))
//
// Log Output Example:
//
//	{
//	  "level": "info",
//	  "msg": "HTTP Request",
//	  "method": "GET",
//	  "path": "/api/v1/urls",
//	  "status": 200,
//	  "duration": "15.2ms",
//	  "user_agent": "Mozilla/5.0...",
//	  "remote_ip": "192.168.1.100"
//	}
func loggingMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now().UTC()

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

// corsMiddleware handles Cross-Origin Resource Sharing
//
// Purpose: Enables cross-origin requests for web applications.
// This middleware sets appropriate CORS headers to allow browsers
// to make requests from different origins to the API Gateway.
//
// Features:
//   - Allows all origins (*)
//   - Supports common HTTP methods (GET, POST, PUT, DELETE, OPTIONS)
//   - Handles preflight OPTIONS requests
//   - Sets appropriate CORS headers
//
// Example Usage:
//
//	router.Use(corsMiddleware())
//
// Headers Set:
//
//	Access-Control-Allow-Origin: *
//	Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
//	Access-Control-Allow-Headers: Content-Type, Authorization
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

// recoveryMiddleware recovers from panics and returns proper error responses
//
// Purpose: Prevents the application from crashing due to unhandled panics.
// This middleware catches any panics that occur during request processing
// and returns a proper HTTP 500 error response instead of crashing the server.
//
// Features:
//   - Panic recovery and logging
//   - Graceful error response
//   - Request context logging for debugging
//
// Example Usage:
//
//	router.Use(recoveryMiddleware(logger))
//
// Error Response:
//
//	{
//	  "error": "Internal server error"
//	}
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

// authMiddleware handles authentication (placeholder for future implementation)
//
// Purpose: Provides authentication and authorization for API endpoints.
// This middleware will validate JWT tokens, API keys, or other authentication
// mechanisms to ensure only authorized users can access protected endpoints.
//
// Current Status: Placeholder implementation that allows all requests
// Future Implementation: JWT validation, API key checking, role-based access control
//
// Example Usage:
//
//	router.Use(authMiddleware(logger))
//
// Future Features:
//   - JWT token validation
//   - API key authentication
//   - Role-based access control
//   - Rate limiting per user
func authMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement JWT authentication
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// rateLimitMiddleware handles rate limiting (placeholder for future implementation)
//
// Purpose: Prevents API abuse by limiting the number of requests per client.
// This middleware will implement rate limiting based on IP address, user ID,
// or other identifiers to protect the API from excessive usage.
//
// Current Status: Placeholder implementation that allows all requests
// Future Implementation: Token bucket algorithm, Redis-based rate limiting
//
// Example Usage:
//
//	router.Use(rateLimitMiddleware(logger))
//
// Future Features:
//   - Token bucket rate limiting
//   - Redis-based distributed rate limiting
//   - Per-endpoint rate limits
//   - Rate limit headers in responses
func rateLimitMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Implement rate limiting
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// requestIDMiddleware adds a unique request ID to each request
//
// Purpose: Provides request tracing and correlation across distributed systems.
// This middleware generates a unique request ID for each incoming request
// and adds it to the response headers for client-side correlation.
//
// Current Status: Placeholder implementation
// Future Implementation: UUID generation, header injection, context propagation
//
// Example Usage:
//
//	router.Use(requestIDMiddleware(logger))
//
// Future Features:
//   - UUID v4 request ID generation
//   - X-Request-ID header injection
//   - Context propagation for internal services
//   - Correlation with logging and metrics
func requestIDMiddleware(log *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Generate and add request ID
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}
