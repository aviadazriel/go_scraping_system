package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// DataHandler handles data-related HTTP requests
type DataHandler struct {
	logger *logrus.Logger
	// TODO: Add service dependencies
	// dataService domain.DataStorageService
}

// NewDataHandler creates a new data handler
func NewDataHandler(logger *logrus.Logger) *DataHandler {
	return &DataHandler{
		logger: logger,
	}
}

// ListDataResponse represents the response for listing data
type ListDataResponse struct {
	Data  []DataItem `json:"data"`
	Total int64      `json:"total"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}

// DataItem represents a data item in the list
type DataItem struct {
	ID        string                 `json:"id"`
	URLID     string                 `json:"url_id"`
	Schema    string                 `json:"schema"`
	Data      map[string]interface{} `json:"data"`
	ScrapedAt string                 `json:"scraped_at"`
	CreatedAt string                 `json:"created_at"`
}

// ListData handles GET /api/v1/data
func (h *DataHandler) ListData(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	schema := r.URL.Query().Get("schema")
	urlID := r.URL.Query().Get("url_id")

	offset := (page - 1) * limit

	// TODO: Get data from service
	// var data []*domain.ParsedData
	// var err error
	//
	// if schema != "" {
	//     data, err = h.dataService.GetBySchema(r.Context(), schema, limit, offset)
	// } else if urlID != "" {
	//     data, err = h.dataService.GetByURLID(r.Context(), urlID, limit, offset)
	// } else {
	//     data, err = h.dataService.GetAll(r.Context(), limit, offset)
	// }
	//
	// if err != nil {
	//     h.logger.WithError(err).Error("Failed to get data")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := ListDataResponse{
		Data:  []DataItem{},
		Total: 0,
		Page:  page,
		Limit: limit,
	}

	// Use variables to avoid unused variable warnings
	_ = offset
	_ = schema
	_ = urlID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDataByURL handles GET /api/v1/data/{url_id}
func (h *DataHandler) GetDataByURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlID := vars["url_id"]

	if urlID == "" {
		http.Error(w, "URL ID is required", http.StatusBadRequest)
		return
	}

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

	// TODO: Get data by URL from service
	// data, err := h.dataService.GetByURLID(r.Context(), urlID, limit, offset)
	// if err != nil {
	//     h.logger.WithError(err).Error("Failed to get data by URL")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := ListDataResponse{
		Data:  []DataItem{},
		Total: 0,
		Page:  page,
		Limit: limit,
	}

	// Use offset to avoid unused variable warning
	_ = offset

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExportDataRequest represents the request body for exporting data
type ExportDataRequest struct {
	Format string   `json:"format" validate:"required,oneof=json csv xml"`
	URLIDs []string `json:"url_ids,omitempty"`
	Schema string   `json:"schema,omitempty"`
	From   string   `json:"from,omitempty"`
	To     string   `json:"to,omitempty"`
}

// ExportData handles GET /api/v1/data/export
func (h *DataHandler) ExportData(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	urlIDs := r.URL.Query()["url_id"]
	schema := r.URL.Query().Get("schema")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	// TODO: Export data using service
	// data, err := h.dataService.ExportData(r.Context(), format, urlIDs, schema, from, to)
	// if err != nil {
	//     h.logger.WithError(err).Error("Failed to export data")
	//     http.Error(w, "Internal server error", http.StatusInternalServerError)
	//     return
	// }

	// For now, return mock data
	response := map[string]interface{}{
		"format": format,
		"count":  0,
		"data":   []interface{}{},
	}

	// Use variables to avoid unused variable warnings
	_ = urlIDs
	_ = schema
	_ = from
	_ = to

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
