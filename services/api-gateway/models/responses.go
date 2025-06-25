package models

// CreateURLResponse represents the response for a successful URL creation.
// It includes the generated ID and basic status information.
type CreateURLResponse struct {
	ID        string `json:"id"`         // Unique identifier for the created URL
	URL       string `json:"url"`        // The original URL that was registered
	Status    string `json:"status"`     // Current status (pending, active, paused, etc.)
	CreatedAt string `json:"created_at"` // ISO 8601 timestamp of creation
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

// ListDataResponse represents the paginated response for listing scraped data.
// It includes the data array and pagination metadata.
type ListDataResponse struct {
	Data  []DataItem `json:"data"`  // Array of data items
	Total int64      `json:"total"` // Total number of data records
	Page  int        `json:"page"`  // Current page number
	Limit int        `json:"limit"` // Number of items per page
}

// DataItem represents a scraped data record in the list response.
// It contains essential information for displaying scraped data.
type DataItem struct {
	ID        string `json:"id"`         // Unique identifier
	URLID     string `json:"url_id"`     // Associated URL ID
	URL       string `json:"url"`        // The URL that was scraped
	Title     string `json:"title"`      // Extracted title
	Content   string `json:"content"`    // Extracted content
	CreatedAt string `json:"created_at"` // When the data was scraped
}

// URLMetricsResponse represents metrics data for a specific URL.
// It provides comprehensive statistics and time series data for URL performance.
type URLMetricsResponse struct {
	URLID               string                `json:"url_id"`             // URL identifier
	TotalScrapes        int64                 `json:"total_scrapes"`      // Total number of scraping attempts
	SuccessfulScrapes   int64                 `json:"successful_scrapes"` // Number of successful scrapes
	FailedScrapes       int64                 `json:"failed_scrapes"`     // Number of failed scrapes
	SuccessRate         float64               `json:"success_rate"`       // Success rate percentage
	AverageResponseTime float64               `json:"avg_response_time"`  // Average response time in milliseconds
	LastScrapeTime      string                `json:"last_scrape_time"`   // Last scrape timestamp
	TimeSeriesData      []TimeSeriesDataPoint `json:"time_series_data"`   // Historical performance data
}

// TimeSeriesDataPoint represents a single data point in time series metrics.
// This struct is used for tracking performance over time.
type TimeSeriesDataPoint struct {
	Timestamp    string  `json:"timestamp"`     // ISO 8601 timestamp
	ResponseTime float64 `json:"response_time"` // Response time in milliseconds
	StatusCode   int     `json:"status_code"`   // HTTP status code
	Success      bool    `json:"success"`       // Whether the scrape was successful
	DataSize     int64   `json:"data_size"`     // Size of scraped data in bytes
}

// SystemMetricsResponse represents system-wide metrics and health information.
// It provides an overview of the entire scraping system's performance.
type SystemMetricsResponse struct {
	TotalURLs           int64   `json:"total_urls"`        // Total number of registered URLs
	ActiveURLs          int64   `json:"active_urls"`       // Number of active URLs
	PendingURLs         int64   `json:"pending_urls"`      // Number of URLs pending scraping
	FailedURLs          int64   `json:"failed_urls"`       // Number of URLs with recent failures
	TotalScrapes        int64   `json:"total_scrapes"`     // Total scraping attempts across all URLs
	SuccessRate         float64 `json:"success_rate"`      // Overall success rate
	AverageResponseTime float64 `json:"avg_response_time"` // Average response time across all URLs
	QueueSize           int64   `json:"queue_size"`        // Current queue size
	WorkerCount         int     `json:"worker_count"`      // Number of active workers
	SystemUptime        string  `json:"system_uptime"`     // System uptime duration
	LastUpdated         string  `json:"last_updated"`      // Last metrics update timestamp
}

// DeadLetterMessageResponse represents a single dead letter message.
// It contains information about a failed message that couldn't be processed.
type DeadLetterMessageResponse struct {
	ID         string `json:"id"`          // Unique message identifier
	Topic      string `json:"topic"`       // Source topic
	Partition  int32  `json:"partition"`   // Kafka partition
	Offset     int64  `json:"offset"`      // Message offset
	Key        string `json:"key"`         // Message key
	Value      string `json:"value"`       // Message value (truncated)
	Error      string `json:"error"`       // Error message
	RetryCount int    `json:"retry_count"` // Number of retry attempts
	CreatedAt  string `json:"created_at"`  // When the message was created
	FailedAt   string `json:"failed_at"`   // When the message failed
}

// ListDeadLetterMessagesResponse represents the paginated response for dead letter messages.
// It includes the messages array and pagination metadata.
type ListDeadLetterMessagesResponse struct {
	Messages []DeadLetterMessageResponse `json:"messages"` // Array of dead letter messages
	Total    int64                       `json:"total"`    // Total number of dead letter messages
	Page     int                         `json:"page"`     // Current page number
	Limit    int                         `json:"limit"`    // Number of items per page
}

// HealthResponse represents the health check response.
// It provides information about the service's health status.
type HealthResponse struct {
	Status    string            `json:"status"`           // Overall health status (healthy, unhealthy)
	Timestamp string            `json:"timestamp"`        // Health check timestamp
	Uptime    string            `json:"uptime"`           // Service uptime duration
	Version   string            `json:"version"`          // Service version
	Checks    map[string]string `json:"checks,omitempty"` // Individual health checks
}
