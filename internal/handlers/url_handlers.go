package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_scraping_project/internal/domain"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// URLHandler handles URL-related HTTP requests
type URLHandler struct {
	logger *logrus.Logger
	// TODO: Add service dependencies
	// urlService domain.URLService
}

// NewURLHandler creates a new URL handler
func NewURLHandler(logger *logrus.Logger) *URLHandler {
	return &URLHandler{
		logger: logger,
	}
}

// CreateURLRequest represents the request body for creating a URL
type CreateURLRequest struct {
	URL          string               `json:"url" validate:"required,url"`
	Frequency    string               `json:"frequency" validate:"required"`
	ParserConfig *domain.ParserConfig `json:"parser_config,omitempty"`
	UserAgent    string               `json:"user_agent,omitempty"`
	Timeout      int                  `json:"timeout,omitempty"`
	RateLimit    int                  `json:"rate_limit,omitempty"`
	MaxRetries   int                  `json:"max_retries,omitempty"`
}

// CreateURLResponse represents the response for creating a URL
type CreateURLResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// CreateURL handles POST /api/v1/urls
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

// ListURLsResponse represents the response for listing URLs
type ListURLsResponse struct {
	URLs  []URLListItem `json:"urls"`
	Total int64         `json:"total"`
	Page  int           `json:"page"`
	Limit int           `json:"limit"`
}

// URLListItem represents a URL in the list
type URLListItem struct {
	ID            string  `json:"id"`
	URL           string  `json:"url"`
	Frequency     string  `json:"frequency"`
	Status        string  `json:"status"`
	LastScrapedAt *string `json:"last_scraped_at,omitempty"`
	NextScrapeAt  *string `json:"next_scrape_at,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// ListURLs handles GET /api/v1/urls
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

// UpdateURLRequest represents the request body for updating a URL
type UpdateURLRequest struct {
	Frequency    string               `json:"frequency,omitempty"`
	ParserConfig *domain.ParserConfig `json:"parser_config,omitempty"`
	UserAgent    string               `json:"user_agent,omitempty"`
	Timeout      int                  `json:"timeout,omitempty"`
	RateLimit    int                  `json:"rate_limit,omitempty"`
	MaxRetries   int                  `json:"max_retries,omitempty"`
}

// UpdateURL handles PUT /api/v1/urls/{id}
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
