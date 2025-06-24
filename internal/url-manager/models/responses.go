package models

// This file is kept for potential future use but is not needed for the background service.
// The URL Manager service now focuses purely on scheduling and task distribution.

import (
	"time"

	"github.com/google/uuid"
)

// TriggerScrapeResponse represents the response for a manual scrape trigger
type TriggerScrapeResponse struct {
	TaskID    uuid.UUID `json:"task_id"`
	URLID     uuid.UUID `json:"url_id"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ScheduleURLResponse represents the response for scheduling a URL
type ScheduleURLResponse struct {
	URLID        uuid.UUID `json:"url_id"`
	URL          string    `json:"url"`
	NextScrapeAt time.Time `json:"next_scrape_at"`
	Status       string    `json:"status"`
}

// BulkTriggerResponse represents the response for bulk trigger operations
type BulkTriggerResponse struct {
	TotalURLs    int                     `json:"total_urls"`
	Triggered    int                     `json:"triggered"`
	Failed       int                     `json:"failed"`
	Tasks        []TriggerScrapeResponse `json:"tasks"`
	FailedURLIDs []uuid.UUID             `json:"failed_url_ids,omitempty"`
}

// ScheduledURLResponse represents a URL scheduled for scraping
type ScheduledURLResponse struct {
	URLID        uuid.UUID `json:"url_id"`
	URL          string    `json:"url"`
	Frequency    string    `json:"frequency"`
	NextScrapeAt time.Time `json:"next_scrape_at"`
	Status       string    `json:"status"`
	RetryCount   int       `json:"retry_count"`
	MaxRetries   int       `json:"max_retries"`
}

// GetScheduledURLsResponse represents the response for getting scheduled URLs
type GetScheduledURLsResponse struct {
	URLs       []ScheduledURLResponse `json:"urls"`
	Total      int                    `json:"total"`
	From       time.Time              `json:"from"`
	To         time.Time              `json:"to"`
	NextScrape *time.Time             `json:"next_scrape,omitempty"`
}

// HealthResponse represents the health status of the URL Manager service
type HealthResponse struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp time.Time         `json:"timestamp"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]string `json:"checks"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
