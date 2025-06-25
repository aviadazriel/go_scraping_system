package database

import (
	"context"
	"database/sql"
	"time"

	"go_scraping_project/shared/models"
)

// Repository defines the interface for database operations
type Repository interface {
	// URL operations
	CreateURL(ctx context.Context, url *models.URL) error
	GetURLByID(ctx context.Context, id string) (*models.URL, error)
	GetAllURLs(ctx context.Context) ([]*models.URL, error)
	UpdateURL(ctx context.Context, url *models.URL) error
	DeleteURL(ctx context.Context, id string) error
	GetURLsForScraping(ctx context.Context) ([]*models.URL, error)
	UpdateURLNextScrapeTime(ctx context.Context, id string, nextScrapeAt *time.Time) error
	UpdateURLLastScrapedTime(ctx context.Context, id string, lastScrapedAt *time.Time) error
	IncrementRetryCount(ctx context.Context, id string) error
	ResetRetryCount(ctx context.Context, id string) error

	// Scraping task operations
	CreateScrapingTask(ctx context.Context, task *models.ScrapingTask) error
	GetScrapingTaskByID(ctx context.Context, id string) (*models.ScrapingTask, error)
	GetAllScrapingTasks(ctx context.Context) ([]*models.ScrapingTask, error)
	UpdateScrapingTask(ctx context.Context, task *models.ScrapingTask) error
	DeleteScrapingTask(ctx context.Context, id string) error

	// Scraped data operations
	CreateScrapedData(ctx context.Context, data *models.ScrapedData) error
	GetScrapedDataByID(ctx context.Context, id string) (*models.ScrapedData, error)
	GetScrapedDataByURLID(ctx context.Context, urlID string) ([]*models.ScrapedData, error)
	GetAllScrapedData(ctx context.Context) ([]*models.ScrapedData, error)
	DeleteScrapedData(ctx context.Context, id string) error

	// Parsed data operations
	CreateParsedData(ctx context.Context, data *models.ParsedData) error
	GetParsedDataByID(ctx context.Context, id string) (*models.ParsedData, error)
	GetParsedDataByURLID(ctx context.Context, urlID string) ([]*models.ParsedData, error)
	GetAllParsedData(ctx context.Context) ([]*models.ParsedData, error)
	DeleteParsedData(ctx context.Context, id string) error
}

// BaseRepository provides common database operations
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// GetDB returns the underlying database connection
func (r *BaseRepository) GetDB() *sql.DB {
	return r.db
}

// Close closes the database connection
func (r *BaseRepository) Close() error {
	return r.db.Close()
}
