package repositories

import (
	"context"
	"time"

	"go_scraping_project/internal/database"

	"github.com/google/uuid"
)

// URLRepository defines the interface for URL data operations
type URLRepository interface {
	// GetURLByID retrieves a URL by its ID
	GetURLByID(ctx context.Context, id uuid.UUID) (*database.Url, error)

	// GetURLsScheduledForScraping retrieves URLs that are scheduled for scraping within a time range
	GetURLsScheduledForScraping(ctx context.Context, from, to time.Time, limit int32) ([]database.Url, error)

	// GetURLsByStatus retrieves URLs by their status
	GetURLsByStatus(ctx context.Context, status string, limit, offset int32) ([]database.Url, error)

	// UpdateURLStatus updates the status of a URL
	UpdateURLStatus(ctx context.Context, id uuid.UUID, status string) error

	// UpdateNextScrapeTime updates the next scrape time for a URL
	UpdateNextScrapeTime(ctx context.Context, id uuid.UUID, nextScrapeAt time.Time) error

	// UpdateLastScrapedTime updates the last scraped time for a URL
	UpdateLastScrapedTime(ctx context.Context, id uuid.UUID, lastScrapedAt time.Time) error

	// IncrementRetryCount increments the retry count for a URL
	IncrementRetryCount(ctx context.Context, id uuid.UUID) error

	// ResetRetryCount resets the retry count for a URL
	ResetRetryCount(ctx context.Context, id uuid.UUID) error

	// GetURLsForImmediateScraping retrieves URLs that should be scraped immediately
	GetURLsForImmediateScraping(ctx context.Context, limit int32) ([]database.Url, error)

	// CountURLsByStatus counts URLs by their status
	CountURLsByStatus(ctx context.Context, status string) (int64, error)

	// GetURLsByIDs retrieves multiple URLs by their IDs
	GetURLsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.Url, error)
}
