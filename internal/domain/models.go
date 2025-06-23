package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// URLStatus represents the status of a URL scraping task
type URLStatus string

const (
	URLStatusPending    URLStatus = "pending"
	URLStatusInProgress URLStatus = "in_progress"
	URLStatusCompleted  URLStatus = "completed"
	URLStatusFailed     URLStatus = "failed"
	URLStatusRetry      URLStatus = "retry"
)

// URL represents a URL to be scraped
type URL struct {
	ID              uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	URL             string          `json:"url" gorm:"not null;uniqueIndex"`
	Frequency       string          `json:"frequency" gorm:"not null"` // e.g., "1h", "24h", "1w"
	LastScrapedAt   *time.Time      `json:"last_scraped_at"`
	NextScrapeAt    *time.Time      `json:"next_scrape_at"`
	Status          URLStatus       `json:"status" gorm:"not null;default:'pending'"`
	RetryCount      int             `json:"retry_count" gorm:"default:0"`
	MaxRetries      int             `json:"max_retries" gorm:"default:3"`
	ParserConfig    json.RawMessage `json:"parser_config" gorm:"type:jsonb"`
	UserAgent       string          `json:"user_agent"`
	Timeout         int             `json:"timeout" gorm:"default:30"` // seconds
	RateLimit       int             `json:"rate_limit" gorm:"default:1"` // requests per second
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"index"`
}

// ScrapingTask represents a task to scrape a URL
type ScrapingTask struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	URLID     uuid.UUID `json:"url_id" gorm:"type:uuid;not null"`
	URL       string    `json:"url" gorm:"not null"`
	Status    URLStatus `json:"status" gorm:"not null;default:'pending'"`
	Attempt   int       `json:"attempt" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScrapedData represents the raw HTML data scraped from a URL
type ScrapedData struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	URLID       uuid.UUID `json:"url_id" gorm:"type:uuid;not null"`
	TaskID      uuid.UUID `json:"task_id" gorm:"type:uuid;not null"`
	HTMLContent string    `json:"html_content" gorm:"type:text"`
	FilePath    string    `json:"file_path" gorm:"not null"`
	FileSize    int64     `json:"file_size"`
	StatusCode  int       `json:"status_code"`
	Headers     json.RawMessage `json:"headers" gorm:"type:jsonb"`
	ScrapedAt   time.Time `json:"scraped_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// ParsedData represents the structured data extracted from HTML
type ParsedData struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	URLID     uuid.UUID       `json:"url_id" gorm:"type:uuid;not null"`
	ScrapedDataID uuid.UUID   `json:"scraped_data_id" gorm:"type:uuid;not null"`
	Data      json.RawMessage `json:"data" gorm:"type:jsonb"`
	Schema    string          `json:"schema" gorm:"not null"` // e.g., "article", "product", "news"
	CreatedAt time.Time       `json:"created_at"`
}

// ParserConfig represents the configuration for parsing HTML
type ParserConfig struct {
	TitleSelector       string            `json:"title_selector,omitempty"`
	ContentSelector     string            `json:"content_selector,omitempty"`
	AuthorSelector      string            `json:"author_selector,omitempty"`
	DateSelector        string            `json:"date_selector,omitempty"`
	ImageSelector       string            `json:"image_selector,omitempty"`
	PriceSelector       string            `json:"price_selector,omitempty"`
	CustomSelectors     map[string]string `json:"custom_selectors,omitempty"`
	ExtractMetadata     bool              `json:"extract_metadata,omitempty"`
	ExtractLinks        bool              `json:"extract_links,omitempty"`
	ExtractImages       bool              `json:"extract_images,omitempty"`
	RemoveScripts       bool              `json:"remove_scripts,omitempty"`
	RemoveStyles        bool              `json:"remove_styles,omitempty"`
	CleanHTML           bool              `json:"clean_html,omitempty"`
}

// DeadLetterMessage represents a failed message that needs to be processed
type DeadLetterMessage struct {
	ID           uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Topic        string          `json:"topic" gorm:"not null"`
	Partition    int32           `json:"partition"`
	Offset       int64           `json:"offset"`
	Key          []byte          `json:"key"`
	Value        []byte          `json:"value"`
	Error        string          `json:"error" gorm:"type:text"`
	RetryCount   int             `json:"retry_count" gorm:"default:0"`
	MaxRetries   int             `json:"max_retries" gorm:"default:3"`
	NextRetryAt  *time.Time      `json:"next_retry_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ScrapingMetrics represents metrics for scraping operations
type ScrapingMetrics struct {
	ID                uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	URLID             uuid.UUID `json:"url_id" gorm:"type:uuid;not null"`
	ResponseTime      int64     `json:"response_time"` // milliseconds
	StatusCode        int       `json:"status_code"`
	ContentLength     int64     `json:"content_length"`
	Success           bool      `json:"success"`
	Error             string    `json:"error"`
	CreatedAt         time.Time `json:"created_at"`
}

// TableName specifies the table name for GORM
func (URL) TableName() string {
	return "urls"
}

func (ScrapingTask) TableName() string {
	return "scraping_tasks"
}

func (ScrapedData) TableName() string {
	return "scraped_data"
}

func (ParsedData) TableName() string {
	return "parsed_data"
}

func (DeadLetterMessage) TableName() string {
	return "dead_letter_messages"
}

func (ScrapingMetrics) TableName() string {
	return "scraping_metrics"
} 