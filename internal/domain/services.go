package domain

import (
	"context"
	"time"
)

// URLService defines the interface for URL management business logic
type URLService interface {
	CreateURL(ctx context.Context, url *URL) error
	GetURL(ctx context.Context, id string) (*URL, error)
	GetAllURLs(ctx context.Context, limit, offset int) ([]*URL, error)
	UpdateURL(ctx context.Context, url *URL) error
	DeleteURL(ctx context.Context, id string) error
	ScheduleScraping(ctx context.Context, id string) error
	GetURLsDueForScraping(ctx context.Context) ([]*URL, error)
	UpdateScrapingStatus(ctx context.Context, id string, status URLStatus) error
	HandleScrapingFailure(ctx context.Context, id string, error string) error
	HandleScrapingSuccess(ctx context.Context, id string) error
}

// ScrapingService defines the interface for scraping business logic
type ScrapingService interface {
	CreateScrapingTask(ctx context.Context, url *URL) (*ScrapingTask, error)
	ProcessScrapingTask(ctx context.Context, task *ScrapingTask) (*ScrapedData, error)
	SaveScrapedData(ctx context.Context, data *ScrapedData) error
	HandleScrapingError(ctx context.Context, task *ScrapingTask, err error) error
	RetryFailedTask(ctx context.Context, taskID string) error
	GetTaskStatus(ctx context.Context, taskID string) (*ScrapingTask, error)
}

// ParsingService defines the interface for HTML parsing business logic
type ParsingService interface {
	ParseHTML(ctx context.Context, scrapedData *ScrapedData, config *ParserConfig) (*ParsedData, error)
	SaveParsedData(ctx context.Context, data *ParsedData) error
	GetParsedData(ctx context.Context, urlID string, limit, offset int) ([]*ParsedData, error)
	GetLatestParsedData(ctx context.Context, urlID string) (*ParsedData, error)
	ValidateParsedData(ctx context.Context, data *ParsedData) error
}

// DeadLetterService defines the interface for dead letter queue management
type DeadLetterService interface {
	ProcessDeadLetterMessage(ctx context.Context, message *DeadLetterMessage) error
	RetryMessage(ctx context.Context, messageID string) error
	GetRetryableMessages(ctx context.Context) ([]*DeadLetterMessage, error)
	DeleteMessage(ctx context.Context, messageID string) error
	HandleRetryFailure(ctx context.Context, message *DeadLetterMessage, err error) error
}

// MetricsService defines the interface for metrics collection and analysis
type MetricsService interface {
	RecordScrapingMetrics(ctx context.Context, metrics *ScrapingMetrics) error
	GetURLMetrics(ctx context.Context, urlID string, days int) (*URLMetrics, error)
	GetSystemMetrics(ctx context.Context) (*SystemMetrics, error)
	CleanupOldMetrics(ctx context.Context, days int) error
}

// SchedulerService defines the interface for URL scheduling
type SchedulerService interface {
	ScheduleURLs(ctx context.Context) error
	CalculateNextScrapeTime(ctx context.Context, frequency string, lastScraped time.Time) (time.Time, error)
	GetOverdueURLs(ctx context.Context) ([]*URL, error)
	UpdateSchedule(ctx context.Context, urlID string, frequency string) error
}

// URLMetrics represents aggregated metrics for a URL
type URLMetrics struct {
	URLID              string  `json:"url_id"`
	TotalScrapes       int64   `json:"total_scrapes"`
	SuccessfulScrapes  int64   `json:"successful_scrapes"`
	FailedScrapes      int64   `json:"failed_scrapes"`
	SuccessRate        float64 `json:"success_rate"`
	AverageResponseTime float64 `json:"average_response_time"`
	LastScrapedAt      *time.Time `json:"last_scraped_at"`
	NextScrapeAt       *time.Time `json:"next_scrape_at"`
}

// SystemMetrics represents system-wide metrics
type SystemMetrics struct {
	TotalURLs          int64   `json:"total_urls"`
	ActiveURLs         int64   `json:"active_urls"`
	PendingTasks       int64   `json:"pending_tasks"`
	FailedTasks        int64   `json:"failed_tasks"`
	DeadLetterMessages int64   `json:"dead_letter_messages"`
	AverageResponseTime float64 `json:"average_response_time"`
	SuccessRate        float64 `json:"success_rate"`
} 