package domain

import (
	"context"
	"time"
)

// URLRepository defines the interface for URL data access
type URLRepository interface {
	Create(ctx context.Context, url *URL) error
	GetByID(ctx context.Context, id string) (*URL, error)
	GetByURL(ctx context.Context, url string) (*URL, error)
	GetAll(ctx context.Context, limit, offset int) ([]*URL, error)
	GetPendingURLs(ctx context.Context) ([]*URL, error)
	GetURLsDueForScraping(ctx context.Context) ([]*URL, error)
	Update(ctx context.Context, url *URL) error
	UpdateStatus(ctx context.Context, id string, status URLStatus) error
	UpdateLastScraped(ctx context.Context, id string, lastScraped time.Time) error
	UpdateNextScrape(ctx context.Context, id string, nextScrape time.Time) error
	IncrementRetryCount(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// ScrapingTaskRepository defines the interface for scraping task data access
type ScrapingTaskRepository interface {
	Create(ctx context.Context, task *ScrapingTask) error
	GetByID(ctx context.Context, id string) (*ScrapingTask, error)
	GetByURLID(ctx context.Context, urlID string) ([]*ScrapingTask, error)
	GetPendingTasks(ctx context.Context, limit int) ([]*ScrapingTask, error)
	UpdateStatus(ctx context.Context, id string, status URLStatus) error
	IncrementAttempt(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// ScrapedDataRepository defines the interface for scraped data access
type ScrapedDataRepository interface {
	Create(ctx context.Context, data *ScrapedData) error
	GetByID(ctx context.Context, id string) (*ScrapedData, error)
	GetByURLID(ctx context.Context, urlID string, limit, offset int) ([]*ScrapedData, error)
	GetByTaskID(ctx context.Context, taskID string) (*ScrapedData, error)
	GetLatestByURLID(ctx context.Context, urlID string) (*ScrapedData, error)
	Update(ctx context.Context, data *ScrapedData) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// ParsedDataRepository defines the interface for parsed data access
type ParsedDataRepository interface {
	Create(ctx context.Context, data *ParsedData) error
	GetByID(ctx context.Context, id string) (*ParsedData, error)
	GetByURLID(ctx context.Context, urlID string, limit, offset int) ([]*ParsedData, error)
	GetByScrapedDataID(ctx context.Context, scrapedDataID string) (*ParsedData, error)
	GetLatestByURLID(ctx context.Context, urlID string) (*ParsedData, error)
	GetBySchema(ctx context.Context, schema string, limit, offset int) ([]*ParsedData, error)
	Update(ctx context.Context, data *ParsedData) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// DeadLetterRepository defines the interface for dead letter message access
type DeadLetterRepository interface {
	Create(ctx context.Context, message *DeadLetterMessage) error
	GetByID(ctx context.Context, id string) (*DeadLetterMessage, error)
	GetRetryableMessages(ctx context.Context) ([]*DeadLetterMessage, error)
	GetByTopic(ctx context.Context, topic string, limit, offset int) ([]*DeadLetterMessage, error)
	UpdateRetryCount(ctx context.Context, id string, retryCount int) error
	UpdateNextRetry(ctx context.Context, id string, nextRetry time.Time) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}

// MetricsRepository defines the interface for scraping metrics access
type MetricsRepository interface {
	Create(ctx context.Context, metrics *ScrapingMetrics) error
	GetByURLID(ctx context.Context, urlID string, limit, offset int) ([]*ScrapingMetrics, error)
	GetSuccessRate(ctx context.Context, urlID string, days int) (float64, error)
	GetAverageResponseTime(ctx context.Context, urlID string, days int) (float64, error)
	GetErrorCount(ctx context.Context, urlID string, days int) (int64, error)
	DeleteOldMetrics(ctx context.Context, days int) error
	Count(ctx context.Context) (int64, error)
} 