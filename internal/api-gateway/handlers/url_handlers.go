package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_scraping_project/internal/domain"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// URLHandler handles URL-related HTTP requests for the web scraping system.
// It provides endpoints for managing URLs that need to be scraped, including
// creation, listing, updating, deletion, and status monitoring.
type URLHandler struct {
	logger *logrus.Logger
	// TODO: Add service dependencies
	// urlService domain.URLService
}

// NewURLHandler creates a new URL handler with the provided logger.
// This function initializes the handler with necessary dependencies.
func NewURLHandler(logger *logrus.Logger) *URLHandler {
	return &URLHandler{
		logger: logger,
	}
}

// CreateURLRequest represents the request body for creating a new URL to be scraped.
// All fields are validated before processing to ensure data integrity.
type CreateURLRequest struct {
	URL          string               `json:"url" validate:"required,url"`   // The URL to be scraped (required)
	Frequency    string               `json:"frequency" validate:"required"` // Scraping frequency (e.g., "1h", "30m", "1d")
	ParserConfig *domain.ParserConfig `json:"parser_config,omitempty"`       // Configuration for parsing scraped content
	UserAgent    string               `json:"user_agent,omitempty"`          // Custom user agent for HTTP requests
	Timeout      int                  `json:"timeout,omitempty"`             // Request timeout in seconds
	RateLimit    int                  `json:"rate_limit,omitempty"`          // Requests per minute limit
	MaxRetries   int                  `json:"max_retries,omitempty"`         // Maximum number of retry attempts
}

// CreateURLResponse represents the response for a successful URL creation.
// It includes the generated ID and basic status information.
type CreateURLResponse struct {
	ID        string `json:"id"`         // Unique identifier for the created URL
	URL       string `json:"url"`        // The original URL that was registered
	Status    string `json:"status"`     // Current status (pending, active, paused, etc.)
	CreatedAt string `json:"created_at"` // ISO 8601 timestamp of creation
}

// CreateURL handles POST /api/v1/urls
//
// Purpose: Registers a new URL to be scraped with the specified configuration.
// This endpoint validates the input, creates a new URL record, and schedules
// it for scraping according to the provided frequency.
//
// Request Body: CreateURLRequest
// Response: CreateURLResponse (201 Created) or error (400/500)
//
// Example Usage:
//
//	POST /api/v1/urls
//	{
//	  "url": "https://example.com",
//	  "frequency": "1h",
//	  "parser_config": {
//	    "selectors": {"title": "h1", "content": ".content"}
//	  }
//	}
func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req CreateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Add validation
	// if err := validate.Struct(req); err != nil {
	//     h.logger.WithError(err).Error("Validation failed")
	//     http.Error(w, "Validation failed", http.StatusBadRequest)
	//     return
	// }

	// TODO: Create URL using service
	// url := &domain.URL{
	//     URL:          req.URL,
	//     Frequency:    req.Frequency,
	//     ParserConfig: req.ParserConfig,
	//     UserAgent:    req.UserAgent,
	//     Timeout:      req.Timeout,
	//     RateLimit:    req.RateLimit,
	//     MaxRetries:   req.MaxRetries,
	// }
	//
	// if err := h.urlService.CreateURL(r.Context(), url); err != nil {
	//     h.logger.WithError(err).Error("Failed to create URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return a mock response
	response := CreateURLResponse{
		ID:        "mock-id",
		URL:       req.URL,
		Status:    "pending",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListURLsResponse represents the paginated response for listing URLs.
// It includes the URLs array and pagination metadata.
type ListURLsResponse struct {
	URLs  []URLListItem `json:"urls"`  // Array of URL items
	Total int64         `json:"total"` // Total number of URLs (for pagination)
	Page  int           `json:"page"`  // Current page number
	Limit int           `json:"limit"` // Number of items per page
}

// URLListItem represents a URL in the list response.
// It contains essential information for displaying URLs in a list view.
type URLListItem struct {
	ID            string  `json:"id"`                        // Unique identifier
	URL           string  `json:"url"`                       // The URL being scraped
	Frequency     string  `json:"frequency"`                 // Scraping frequency
	Status        string  `json:"status"`                    // Current status
	LastScrapedAt *string `json:"last_scraped_at,omitempty"` // Last successful scrape time
	NextScrapeAt  *string `json:"next_scrape_at,omitempty"`  // Next scheduled scrape time
	CreatedAt     string  `json:"created_at"`                // Creation timestamp
}

// ListURLs handles GET /api/v1/urls
//
// Purpose: Retrieves a paginated list of all registered URLs for scraping.
// This endpoint supports pagination and can be used for dashboard displays
// or administrative interfaces.
//
// Query Parameters:
//   - page: Page number (default: 1)
//   - limit: Items per page, max 100 (default: 20)
//
// Response: ListURLsResponse (200 OK) or error (500)
//
// Example Usage:
//
//	GET /api/v1/urls?page=1&limit=20
func (h *URLHandler) ListURLs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// TODO: Get URLs from service
	// urls, err := h.urlService.GetAllURLs(r.Context(), limit, offset)
	// if err != nil {
	//     h.logger.WithError(err).Error("Failed to get URLs")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := ListURLsResponse{
		URLs:  []URLListItem{},
		Total: 0,
		Page:  page,
		Limit: limit,
	}

	// Use offset to avoid unused variable warning
	_ = offset

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetURL handles GET /api/v1/urls/{id}
//
// Purpose: Retrieves detailed information about a specific URL by its ID.
// This endpoint provides comprehensive information including configuration,
// status, and timing details for a single URL.
//
// Path Parameters:
//   - id: URL identifier (required)
//
// Response: URL details (200 OK) or error (400/404/500)
//
// Example Usage:
//
//	GET /api/v1/urls/url-123
func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Get URL from service
	// url, err := h.urlService.GetURL(r.Context(), id)
	// if err != nil {
	//     if errors.Is(err, domain.ErrURLNotFound) {
	//         http.Error(w, "URL not found", http.StatusNotFound)
	//         return
	//     }
	//     h.logger.WithError(err).Error("Failed to get URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := map[string]interface{}{
		"id":         id,
		"url":        "https://example.com",
		"frequency":  "1h",
		"status":     "pending",
		"created_at": "2024-01-01T00:00:00Z",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateURLRequest represents the request body for updating an existing URL.
// All fields are optional, allowing partial updates of URL configuration.
type UpdateURLRequest struct {
	Frequency    string               `json:"frequency,omitempty"`     // New scraping frequency
	ParserConfig *domain.ParserConfig `json:"parser_config,omitempty"` // Updated parser configuration
	UserAgent    string               `json:"user_agent,omitempty"`    // New user agent
	Timeout      int                  `json:"timeout,omitempty"`       // New timeout value
	RateLimit    int                  `json:"rate_limit,omitempty"`    // New rate limit
	MaxRetries   int                  `json:"max_retries,omitempty"`   // New max retries
}

// UpdateURL handles PUT /api/v1/urls/{id}
//
// Purpose: Updates configuration for an existing URL. This endpoint supports
// partial updates, allowing clients to modify only specific fields without
// providing the complete URL configuration.
//
// Path Parameters:
//   - id: URL identifier (required)
//
// Request Body: UpdateURLRequest (all fields optional)
// Response: Success message (200 OK) or error (400/404/500)
//
// Example Usage:
//
//	PUT /api/v1/urls/url-123
//	{
//	  "frequency": "2h",
//	  "timeout": 45
//	}
func (h *URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

	var req UpdateURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.WithError(err).Error("Failed to decode request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: Update URL using service
	// url, err := h.urlService.GetURL(r.Context(), id)
	// if err != nil {
	//     if errors.Is(err, domain.ErrURLNotFound) {
	//         http.Error(w, "URL not found", http.StatusNotFound)
	//         return
	//     }
	//     h.logger.WithError(err).Error("Failed to get URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }
	//
	// // Update fields
	// if req.Frequency != "" {
	//     url.Frequency = req.Frequency
	// }
	// if req.ParserConfig != nil {
	//     url.ParserConfig = req.ParserConfig
	// }
	// if req.UserAgent != "" {
	//     url.UserAgent = req.UserAgent
	// }
	// if req.Timeout > 0 {
	//     url.Timeout = req.Timeout
	// }
	// if req.RateLimit > 0 {
	//     url.RateLimit = req.RateLimit
	// }
	// if req.MaxRetries > 0 {
	//     url.MaxRetries = req.MaxRetries
	// }
	//
	// if err := h.urlService.UpdateURL(r.Context(), url); err != nil {
	//     h.logger.WithError(err).Error("Failed to update URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "URL updated successfully"})
}

// DeleteURL handles DELETE /api/v1/urls/{id}
//
// Purpose: Removes a URL from the scraping schedule. This operation is
// irreversible and will stop all future scraping attempts for this URL.
// Existing scraped data is preserved unless explicitly configured otherwise.
//
// Path Parameters:
//   - id: URL identifier (required)
//
// Response: Success message (200 OK) or error (400/404/500)
//
// Example Usage:
//
//	DELETE /api/v1/urls/url-123
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Delete URL using service
	// if err := h.urlService.DeleteURL(r.Context(), id); err != nil {
	//     if errors.Is(err, domain.ErrURLNotFound) {
	//         http.Error(w, "URL not found", http.StatusNotFound)
	//         return
	//     }
	//     h.logger.WithError(err).Error("Failed to delete URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "URL deleted successfully"})
}

// TriggerScrape handles POST /api/v1/urls/{id}/scrape
//
// Purpose: Manually triggers scraping for a specific URL, bypassing the
// normal schedule. This is useful for immediate data collection or
// testing purposes. The scraping will be queued and processed as soon
// as a worker becomes available.
//
// Path Parameters:
//   - id: URL identifier (required)
//
// Response: Success message (200 OK) or error (400/404/500)
//
// Example Usage:
//
//	POST /api/v1/urls/url-123/scrape
func (h *URLHandler) TriggerScrape(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Trigger immediate scrape using service
	// if err := h.urlService.ScheduleScraping(r.Context(), id); err != nil {
	//     if errors.Is(err, domain.ErrURLNotFound) {
	//         http.Error(w, "URL not found", http.StatusNotFound)
	//         return
	//     }
	//     h.logger.WithError(err).Error("Failed to trigger scrape")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Scrape triggered successfully"})
}

// GetURLStatus handles GET /api/v1/urls/{id}/status
//
// Purpose: Retrieves current status and scheduling information for a URL.
// This endpoint provides real-time information about the URL's scraping
// status, including last scrape time, next scheduled scrape, and retry
// information.
//
// Path Parameters:
//   - id: URL identifier (required)
//
// Response: URL status details (200 OK) or error (400/404/500)
//
// Example Usage:
//
//	GET /api/v1/urls/url-123/status
func (h *URLHandler) GetURLStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

	// TODO: Get URL status using service
	// url, err := h.urlService.GetURL(r.Context(), id)
	// if err != nil {
	//     if errors.Is(err, domain.ErrURLNotFound) {
	//         http.Error(w, "URL not found", http.StatusNotFound)
	//         return
	//     }
	//     h.logger.WithError(err).Error("Failed to get URL status")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := map[string]interface{}{
		"id":              id,
		"status":          "pending",
		"last_scraped_at": nil,
		"next_scrape_at":  "2024-01-01T01:00:00Z",
		"retry_count":     0,
		"max_retries":     3,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
