package repositories

import (
	"context"
	"database/sql"
	"time"

	"go_scraping_project/internal/database"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// URLRepositoryImpl implements the URLRepository interface using sqlc-generated queries
type URLRepositoryImpl struct {
	db     database.Querier
	logger *logrus.Logger
}

// NewURLRepository creates a new URL repository instance
func NewURLRepository(db database.Querier, logger *logrus.Logger) URLRepository {
	return &URLRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// GetURLByID retrieves a URL by its ID
func (r *URLRepositoryImpl) GetURLByID(ctx context.Context, id uuid.UUID) (*database.Url, error) {
	url, err := r.db.GetURLByID(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("url_id", id).Error("Failed to get URL by ID")
		return nil, err
	}
	return &url, nil
}

// GetURLsScheduledForScraping retrieves URLs that are scheduled for scraping within a time range
func (r *URLRepositoryImpl) GetURLsScheduledForScraping(ctx context.Context, from, to time.Time, limit int32) ([]database.Url, error) {
	urls, err := r.db.GetURLsScheduledForScraping(ctx, database.GetURLsScheduledForScrapingParams{
		NextScrapeAt:   sql.NullTime{Time: from, Valid: true},
		NextScrapeAt_2: sql.NullTime{Time: to, Valid: true},
		Limit:          limit,
	})
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"from":  from,
			"to":    to,
			"limit": limit,
		}).Error("Failed to get URLs scheduled for scraping")
		return nil, err
	}
	return urls, nil
}

// GetURLsByStatus retrieves URLs by their status
func (r *URLRepositoryImpl) GetURLsByStatus(ctx context.Context, status string, limit, offset int32) ([]database.Url, error) {
	urls, err := r.db.GetURLsByStatus(ctx, database.GetURLsByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"status": status,
			"limit":  limit,
			"offset": offset,
		}).Error("Failed to get URLs by status")
		return nil, err
	}
	return urls, nil
}

// UpdateURLStatus updates the status of a URL
func (r *URLRepositoryImpl) UpdateURLStatus(ctx context.Context, id uuid.UUID, status string) error {
	err := r.db.UpdateURLStatus(ctx, database.UpdateURLStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"url_id": id,
			"status": status,
		}).Error("Failed to update URL status")
		return err
	}
	return nil
}

// UpdateNextScrapeTime updates the next scrape time for a URL
func (r *URLRepositoryImpl) UpdateNextScrapeTime(ctx context.Context, id uuid.UUID, nextScrapeAt time.Time) error {
	err := r.db.UpdateNextScrapeTime(ctx, database.UpdateNextScrapeTimeParams{
		ID:           id,
		NextScrapeAt: sql.NullTime{Time: nextScrapeAt, Valid: true},
	})
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"url_id":         id,
			"next_scrape_at": nextScrapeAt,
		}).Error("Failed to update next scrape time")
		return err
	}
	return nil
}

// UpdateLastScrapedTime updates the last scraped time for a URL
func (r *URLRepositoryImpl) UpdateLastScrapedTime(ctx context.Context, id uuid.UUID, lastScrapedAt time.Time) error {
	err := r.db.UpdateLastScrapedTime(ctx, database.UpdateLastScrapedTimeParams{
		ID:            id,
		LastScrapedAt: sql.NullTime{Time: lastScrapedAt, Valid: true},
	})
	if err != nil {
		r.logger.WithError(err).WithFields(logrus.Fields{
			"url_id":          id,
			"last_scraped_at": lastScrapedAt,
		}).Error("Failed to update last scraped time")
		return err
	}
	return nil
}

// IncrementRetryCount increments the retry count for a URL
func (r *URLRepositoryImpl) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	err := r.db.IncrementRetryCount(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("url_id", id).Error("Failed to increment retry count")
		return err
	}
	return nil
}

// ResetRetryCount resets the retry count for a URL
func (r *URLRepositoryImpl) ResetRetryCount(ctx context.Context, id uuid.UUID) error {
	err := r.db.ResetRetryCount(ctx, id)
	if err != nil {
		r.logger.WithError(err).WithField("url_id", id).Error("Failed to reset retry count")
		return err
	}
	return nil
}

// GetURLsForImmediateScraping retrieves URLs that should be scraped immediately
func (r *URLRepositoryImpl) GetURLsForImmediateScraping(ctx context.Context, limit int32) ([]database.Url, error) {
	urls, err := r.db.GetURLsForImmediateScraping(ctx, database.GetURLsForImmediateScrapingParams{
		NextScrapeAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		Limit:        limit,
	})
	if err != nil {
		r.logger.WithError(err).WithField("limit", limit).Error("Failed to get URLs for immediate scraping")
		return nil, err
	}
	return urls, nil
}

// CountURLsByStatus counts URLs by their status
func (r *URLRepositoryImpl) CountURLsByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.db.CountURLsByStatus(ctx, status)
	if err != nil {
		r.logger.WithError(err).WithField("status", status).Error("Failed to count URLs by status")
		return 0, err
	}
	return count, nil
}

// GetURLsByIDs retrieves multiple URLs by their IDs
func (r *URLRepositoryImpl) GetURLsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.Url, error) {
	urls, err := r.db.GetURLsByIDs(ctx, ids)
	if err != nil {
		r.logger.WithError(err).WithField("url_ids", ids).Error("Failed to get URLs by IDs")
		return nil, err
	}
	return urls, nil
}
