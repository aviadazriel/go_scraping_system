package services

import (
	"context"
	"fmt"
	"time"

	"go_scraping_project/internal/database"
	"go_scraping_project/internal/domain"
	"go_scraping_project/internal/url-manager/models"
	"go_scraping_project/internal/url-manager/repositories"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// URLSchedulerService handles URL scheduling and scraping task creation
type URLSchedulerService struct {
	urlRepo   repositories.URLRepository
	producer  domain.KafkaProducer
	logger    *logrus.Logger
	scheduler *time.Ticker
	stopChan  chan struct{}
}

// NewURLSchedulerService creates a new URL scheduler service
func NewURLSchedulerService(
	urlRepo repositories.URLRepository,
	producer domain.KafkaProducer,
	logger *logrus.Logger,
) *URLSchedulerService {
	return &URLSchedulerService{
		urlRepo:  urlRepo,
		producer: producer,
		logger:   logger,
		stopChan: make(chan struct{}),
	}
}

// Start starts the URL scheduler service
func (s *URLSchedulerService) Start(ctx context.Context) error {
	s.logger.Info("Starting URL Scheduler Service")

	// Start the scheduler ticker (check every 30 seconds)
	s.scheduler = time.NewTicker(30 * time.Second)

	go s.runScheduler(ctx)

	return nil
}

// Stop stops the URL scheduler service
func (s *URLSchedulerService) Stop() error {
	s.logger.Info("Stopping URL Scheduler Service")

	if s.scheduler != nil {
		s.scheduler.Stop()
	}

	close(s.stopChan)
	return nil
}

// runScheduler runs the main scheduling loop
func (s *URLSchedulerService) runScheduler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context cancelled, stopping scheduler")
			return
		case <-s.stopChan:
			s.logger.Info("Stop signal received, stopping scheduler")
			return
		case <-s.scheduler.C:
			if err := s.processScheduledURLs(ctx); err != nil {
				s.logger.WithError(err).Error("Failed to process scheduled URLs")
			}
		}
	}
}

// processScheduledURLs processes URLs that are scheduled for scraping
func (s *URLSchedulerService) processScheduledURLs(ctx context.Context) error {
	// Use UTC for all time calculations
	now := time.Now().UTC()
	from := now.Add(-1 * time.Minute) // Include URLs that were due up to 1 minute ago
	to := now.Add(5 * time.Minute)    // Include URLs due in the next 5 minutes

	urls, err := s.urlRepo.GetURLsScheduledForScraping(ctx, from, to, 100)
	if err != nil {
		return fmt.Errorf("failed to get scheduled URLs: %w", err)
	}

	if len(urls) == 0 {
		return nil
	}

	s.logger.WithField("url_count", len(urls)).Info("Processing scheduled URLs")

	for _, url := range urls {
		if err := s.processURL(ctx, url); err != nil {
			s.logger.WithError(err).WithField("url_id", url.ID).Error("Failed to process URL")
			continue
		}
	}

	return nil
}

// processURL processes a single URL for scraping
func (s *URLSchedulerService) processURL(ctx context.Context, url database.Url) error {
	// Check if URL is actually due (double-check to avoid race conditions)
	if !url.NextScrapeAt.Valid || url.NextScrapeAt.Time.After(time.Now().UTC()) {
		return nil // Not actually due yet
	}

	s.logger.Printf("Processing URL: %s (ID: %s)", url.Url, url.ID)

	// Create scraping task message
	taskID := uuid.New()
	message := domain.ScrapingTaskMessage{
		ID:        taskID,
		URLID:     url.ID,
		URL:       url.Url,
		UserAgent: url.UserAgent,
		Timeout:   url.Timeout,
		CreatedAt: time.Now().UTC(),
	}

	// Send message to Kafka
	if err := s.producer.SendMessage(ctx, "scraping-tasks", taskID.String(), message); err != nil {
		return fmt.Errorf("failed to send scraping task to Kafka: %w", err)
	}

	s.logger.Printf("Sent scraping task to Kafka: %s", taskID)

	// Update URL status and last scraped time
	if err := s.urlRepo.UpdateLastScrapedTime(ctx, url.ID, time.Now().UTC()); err != nil {
		return fmt.Errorf("failed to update last scraped time: %w", err)
	}

	// Calculate next scrape time
	nextScrape, err := models.CalculateNextScrapeTime(url.Frequency, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to calculate next scrape time: %w", err)
	}

	// Update next scrape time
	if err := s.urlRepo.UpdateNextScrapeTime(ctx, url.ID, nextScrape); err != nil {
		return fmt.Errorf("failed to update next scrape time: %w", err)
	}

	s.logger.Printf("Updated URL %s: next scrape at %s", url.ID, nextScrape.UTC().Format(time.RFC3339))

	return nil
}
