package database

import (
	"context"

	"github.com/google/uuid"
)

// Querier defines the interface for database operations
// This interface matches the methods we need from the sqlc-generated Queries type
type Querier interface {
	// URL operations
	GetURLByID(ctx context.Context, id uuid.UUID) (Url, error)
	GetURLsScheduledForScraping(ctx context.Context, arg GetURLsScheduledForScrapingParams) ([]Url, error)
	GetURLsByStatus(ctx context.Context, arg GetURLsByStatusParams) ([]Url, error)
	UpdateURLStatus(ctx context.Context, arg UpdateURLStatusParams) error
	UpdateNextScrapeTime(ctx context.Context, arg UpdateNextScrapeTimeParams) error
	UpdateLastScrapedTime(ctx context.Context, arg UpdateLastScrapedTimeParams) error
	IncrementRetryCount(ctx context.Context, id uuid.UUID) error
	ResetRetryCount(ctx context.Context, id uuid.UUID) error
	GetURLsForImmediateScraping(ctx context.Context, arg GetURLsForImmediateScrapingParams) ([]Url, error)
	CountURLsByStatus(ctx context.Context, status string) (int64, error)
	GetURLsByIDs(ctx context.Context, dollar_1 []uuid.UUID) ([]Url, error)
}
