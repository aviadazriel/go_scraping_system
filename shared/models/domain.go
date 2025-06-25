package models

import (
	"time"

	"github.com/google/uuid"
)

// Domain models shared across all services

// URL represents a URL to be scraped
type URL struct {
	ID            uuid.UUID     `json:"id"`
	URL           string        `json:"url"`
	Frequency     string        `json:"frequency"`
	Status        string        `json:"status"`
	MaxRetries    int           `json:"max_retries"`
	Timeout       int           `json:"timeout"`
	RateLimit     int           `json:"rate_limit"`
	UserAgent     string        `json:"user_agent,omitempty"`
	ParserConfig  *ParserConfig `json:"parser_config,omitempty"`
	NextScrapeAt  *time.Time    `json:"next_scrape_at,omitempty"`
	LastScrapedAt *time.Time    `json:"last_scraped_at,omitempty"`
	RetryCount    int           `json:"retry_count"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// ParserConfig represents configuration for parsing scraped content
type ParserConfig struct {
	Selectors map[string]string `json:"selectors"`       // CSS selectors for different content types
	Rules     []ParseRule       `json:"rules,omitempty"` // Custom parsing rules
}

// ParseRule represents a custom parsing rule
type ParseRule struct {
	Name     string `json:"name"`
	Selector string `json:"selector"`
	Type     string `json:"type"`           // text, attr, html, etc.
	Attr     string `json:"attr,omitempty"` // attribute name for attr type
}

// ScrapingTask represents a task to scrape a URL
type ScrapingTask struct {
	ID         uuid.UUID `json:"id"`
	URLID      uuid.UUID `json:"url_id"`
	URL        string    `json:"url"`
	UserAgent  string    `json:"user_agent"`
	Timeout    int       `json:"timeout"`
	MaxRetries int       `json:"max_retries"`
	CreatedAt  time.Time `json:"created_at"`
}

// ScrapedData represents raw scraped data
type ScrapedData struct {
	ID          uuid.UUID `json:"id"`
	URLID       uuid.UUID `json:"url_id"`
	URL         string    `json:"url"`
	StatusCode  int       `json:"status_code"`
	Content     string    `json:"content"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	Duration    float64   `json:"duration"` // in milliseconds
	CreatedAt   time.Time `json:"created_at"`
}

// ParsedData represents parsed/structured data
type ParsedData struct {
	ID        uuid.UUID              `json:"id"`
	URLID     uuid.UUID              `json:"url_id"`
	URL       string                 `json:"url"`
	Title     string                 `json:"title,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// KafkaMessage represents a generic Kafka message
type KafkaMessage struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Common status values
const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusPaused  = "paused"
	StatusFailed  = "failed"
	StatusSuccess = "success"
)

// Common frequency values
const (
	FrequencyMinute = "1m"
	FrequencyHour   = "1h"
	FrequencyDay    = "1d"
	FrequencyWeek   = "1w"
)
