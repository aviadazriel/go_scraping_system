package models

// This file is kept for potential future use but is not needed for the background service.
// The URL Manager service now focuses purely on scheduling and task distribution.

import (
	"time"

	"github.com/google/uuid"
)

// TriggerScrapeRequest represents a request to manually trigger scraping for a URL
type TriggerScrapeRequest struct {
	URLID uuid.UUID `json:"url_id" validate:"required"`
}

// ScheduleURLRequest represents a request to schedule a URL for scraping
type ScheduleURLRequest struct {
	URLID uuid.UUID `json:"url_id" validate:"required"`
}

// UpdateNextScrapeRequest represents a request to update the next scrape time for a URL
type UpdateNextScrapeRequest struct {
	URLID        uuid.UUID `json:"url_id" validate:"required"`
	NextScrapeAt time.Time `json:"next_scrape_at" validate:"required"`
}

// BulkTriggerRequest represents a request to trigger scraping for multiple URLs
type BulkTriggerRequest struct {
	URLIDs []uuid.UUID `json:"url_ids" validate:"required,min=1,max=100"`
}

// GetScheduledURLsRequest represents a request to get URLs scheduled for scraping
type GetScheduledURLsRequest struct {
	From  time.Time `json:"from" validate:"required"`
	To    time.Time `json:"to" validate:"required"`
	Limit int       `json:"limit" validate:"min=1,max=1000"`
}
