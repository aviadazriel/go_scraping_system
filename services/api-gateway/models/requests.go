package models

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

// ExportDataRequest represents the request body for exporting scraped data.
// This struct defines the parameters for data export operations.
type ExportDataRequest struct {
	Format    string   `json:"format" validate:"required,oneof=json csv xml"` // Export format (json, csv, xml)
	URLIDs    []string `json:"url_ids,omitempty"`                             // Specific URL IDs to export
	StartDate string   `json:"start_date,omitempty"`                          // Start date for data range (ISO 8601)
	EndDate   string   `json:"end_date,omitempty"`                            // End date for data range (ISO 8601)
	Limit     int      `json:"limit,omitempty"`                               // Maximum number of records to export
}

// BulkRetryRequest represents the request body for bulk retry operations.
// This struct defines parameters for retrying multiple failed messages.
type BulkRetryRequest struct {
	MessageIDs []string `json:"message_ids" validate:"required,min=1"` // Array of message IDs to retry
	Topic      string   `json:"topic,omitempty"`                       // Target topic for retry (optional)
}
