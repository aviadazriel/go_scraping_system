package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go_scraping_project/internal/database"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sqlc-dev/pqtype"
)

// ParserConfig represents the configuration for parsing HTML
// This is a simplified version for the API Gateway
type ParserConfig struct {
	TitleSelector   string            `json:"title_selector,omitempty"`
	ContentSelector string            `json:"content_selector,omitempty"`
	AuthorSelector  string            `json:"author_selector,omitempty"`
	DateSelector    string            `json:"date_selector,omitempty"`
	ImageSelector   string            `json:"image_selector,omitempty"`
	PriceSelector   string            `json:"price_selector,omitempty"`
	CustomSelectors map[string]string `json:"custom_selectors,omitempty"`
	ExtractMetadata bool              `json:"extract_metadata,omitempty"`
	ExtractLinks    bool              `json:"extract_links,omitempty"`
	ExtractImages   bool              `json:"extract_images,omitempty"`
	RemoveScripts   bool              `json:"remove_scripts,omitempty"`
	RemoveStyles    bool              `json:"remove_styles,omitempty"`
	CleanHTML       bool              `json:"clean_html,omitempty"`
}

// URLHandler handles URL-related HTTP requests for the web scraping system.
// It provides endpoints for managing URLs that need to be scraped, including
// creation, listing, updating, deletion, and status monitoring.
type URLHandler struct {
	logger *logrus.Logger
	db     *database.Queries // sqlc-generated database queries
}

// NewURLHandler creates a new URL handler with the provided logger and database queries.
// This function initializes the handler with necessary dependencies for URL management.
func NewURLHandler(logger *logrus.Logger, db *database.Queries) *URLHandler {
	return &URLHandler{
		logger: logger,
		db:     db,
	}
}

// CreateURLRequest represents the request body for creating a new URL to be scraped.
// All fields are validated before processing to ensure data integrity.
type CreateURLRequest struct {
	URL          string        `json:"url" validate:"required,url"`   // The URL to be scraped (required)
	Frequency    string        `json:"frequency" validate:"required"` // Scraping frequency (e.g., "1h", "30m", "1d")
	ParserConfig *ParserConfig `json:"parser_config,omitempty"`       // Configuration for parsing scraped content
	UserAgent    string        `json:"user_agent,omitempty"`          // Custom user agent for HTTP requests
	Timeout      int           `json:"timeout,omitempty"`             // Request timeout in seconds
	RateLimit    int           `json:"rate_limit,omitempty"`          // Requests per minute limit
	MaxRetries   int           `json:"max_retries,omitempty"`         // Maximum number of retry attempts
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
// This endpoint validates the input, creates a new URL record in the database,
// and returns the created URL with its generated ID.
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

	// Validate request
	if err := h.validateCreateURLRequest(&req); err != nil {
		h.logger.WithError(err).WithField("url", req.URL).Error("Validation failed")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate next scrape time
	nextScrape, err := h.calculateNextScrapeTime(req.Frequency, time.Now())
	if err != nil {
		h.logger.WithError(err).Error("Failed to calculate next scrape time")
		http.Error(w, "Invalid frequency format", http.StatusBadRequest)
		return
	}

	// Prepare parser config JSON if provided
	var parserConfigJSON pqtype.NullRawMessage
	if req.ParserConfig != nil {
		configBytes, err := json.Marshal(req.ParserConfig)
		if err != nil {
			h.logger.WithError(err).Error("Failed to marshal parser config")
			http.Error(w, "Invalid parser configuration", http.StatusBadRequest)
			return
		}
		parserConfigJSON = pqtype.NullRawMessage{
			RawMessage: configBytes,
			Valid:      true,
		}
	}

	// Prepare user agent
	var userAgent sql.NullString
	if req.UserAgent != "" {
		userAgent = sql.NullString{
			String: req.UserAgent,
			Valid:  true,
		}
	} else {
		userAgent = sql.NullString{
			String: "GoScrapingBot/1.0",
			Valid:  true,
		}
	}

	// Create URL using sqlc-generated database queries
	params := database.CreateURLParams{
		Url:          req.URL,
		Frequency:    req.Frequency,
		Status:       "pending",
		MaxRetries:   int32(h.getDefaultValue(req.MaxRetries, 3)),
		Timeout:      int32(h.getDefaultValue(req.Timeout, 30)),
		RateLimit:    int32(h.getDefaultValue(req.RateLimit, 1)),
		UserAgent:    userAgent,
		ParserConfig: parserConfigJSON,
		NextScrapeAt: sql.NullTime{
			Time:  nextScrape,
			Valid: true,
		},
	}

	createdURL, err := h.db.CreateURL(r.Context(), params)
	if err != nil {
		h.logger.WithError(err).WithField("url", req.URL).Error("Failed to save URL to database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := CreateURLResponse{
		ID:        createdURL.ID.String(),
		URL:       createdURL.Url,
		Status:    createdURL.Status,
		CreatedAt: createdURL.CreatedAt.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// validateCreateURLRequest validates the CreateURLRequest
// This function performs comprehensive validation of the request data
// including URL format, frequency format, and business rule validation.
func (h *URLHandler) validateCreateURLRequest(req *CreateURLRequest) error {
	// Validate URL
	if req.URL == "" {
		return &ValidationError{Field: "url", Message: "URL is required"}
	}

	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return &ValidationError{Field: "url", Message: "Invalid URL format"}
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return &ValidationError{Field: "url", Message: "URL must include scheme and host"}
	}

	// Validate frequency
	if req.Frequency == "" {
		return &ValidationError{Field: "frequency", Message: "Frequency is required"}
	}

	if err := h.validateFrequency(req.Frequency); err != nil {
		return &ValidationError{Field: "frequency", Message: err.Error()}
	}

	// Validate timeout
	if req.Timeout < 0 {
		return &ValidationError{Field: "timeout", Message: "Timeout must be non-negative"}
	}

	if req.Timeout > 300 { // Max 5 minutes
		return &ValidationError{Field: "timeout", Message: "Timeout cannot exceed 300 seconds"}
	}

	// Validate rate limit
	if req.RateLimit < 0 {
		return &ValidationError{Field: "rate_limit", Message: "Rate limit must be non-negative"}
	}

	if req.RateLimit > 1000 { // Max 1000 requests per minute
		return &ValidationError{Field: "rate_limit", Message: "Rate limit cannot exceed 1000 requests per minute"}
	}

	// Validate max retries
	if req.MaxRetries < 0 {
		return &ValidationError{Field: "max_retries", Message: "Max retries must be non-negative"}
	}

	if req.MaxRetries > 10 { // Max 10 retries
		return &ValidationError{Field: "max_retries", Message: "Max retries cannot exceed 10"}
	}

	return nil
}

// validateFrequency validates the frequency string format
// This function ensures the frequency follows the expected format (e.g., "1h", "30m", "1d").
func (h *URLHandler) validateFrequency(frequency string) error {
	if frequency == "" {
		return &ValidationError{Field: "frequency", Message: "Frequency cannot be empty"}
	}

	// Check if frequency ends with a valid unit
	validUnits := []string{"s", "m", "h", "d", "w"}
	hasValidUnit := false

	for _, unit := range validUnits {
		if strings.HasSuffix(frequency, unit) {
			hasValidUnit = true
			break
		}
	}

	if !hasValidUnit {
		return &ValidationError{Field: "frequency", Message: "Frequency must end with a valid unit (s, m, h, d, w)"}
	}

	// Extract numeric part
	numericPart := strings.TrimSuffix(frequency, frequency[len(frequency)-1:])
	if numericPart == "" {
		return &ValidationError{Field: "frequency", Message: "Frequency must include a numeric value"}
	}

	// Parse numeric value
	value, err := strconv.Atoi(numericPart)
	if err != nil {
		return &ValidationError{Field: "frequency", Message: "Frequency must be a valid number"}
	}

	if value <= 0 {
		return &ValidationError{Field: "frequency", Message: "Frequency value must be positive"}
	}

	// Validate minimum frequency (at least 30 seconds)
	if strings.HasSuffix(frequency, "s") && value < 30 {
		return &ValidationError{Field: "frequency", Message: "Minimum frequency is 30 seconds"}
	}

	return nil
}

// getDefaultValue returns the default value if the input is 0, otherwise returns the input
// This helper function provides sensible defaults for optional numeric fields.
func (h *URLHandler) getDefaultValue(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

// ValidationError represents a validation error with field-specific information
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error implements the error interface for ValidationError
func (e *ValidationError) Error() string {
	return e.Message
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
	Frequency    string        `json:"frequency,omitempty"`     // New scraping frequency
	ParserConfig *ParserConfig `json:"parser_config,omitempty"` // Updated parser configuration
	UserAgent    string        `json:"user_agent,omitempty"`    // New user agent
	Timeout      int           `json:"timeout,omitempty"`       // New timeout value
	RateLimit    int           `json:"rate_limit,omitempty"`    // New rate limit
	MaxRetries   int           `json:"max_retries,omitempty"`   // New max retries
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

// calculateNextScrapeTime calculates when the URL should be scraped next
func (h *URLHandler) calculateNextScrapeTime(frequency string, from time.Time) (time.Time, error) {
	duration, err := h.parseFrequency(frequency)
	if err != nil {
		return time.Time{}, err
	}
	return from.Add(duration), nil
}

// parseFrequency parses frequency string into time.Duration
func (h *URLHandler) parseFrequency(frequency string) (time.Duration, error) {
	switch frequency {
	case "30s":
		return 30 * time.Second, nil
	case "1m":
		return 1 * time.Minute, nil
	case "5m":
		return 5 * time.Minute, nil
	case "15m":
		return 15 * time.Minute, nil
	case "30m":
		return 30 * time.Minute, nil
	case "1h":
		return 1 * time.Hour, nil
	case "6h":
		return 6 * time.Hour, nil
	case "12h":
		return 12 * time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	case "1w":
		return 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported frequency: %s", frequency)
	}
}
