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
	now := time.Now()
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
	// Check if the URL is due for scraping
	if !url.NextScrapeAt.Valid || url.NextScrapeAt.Time.After(time.Now()) {
		return nil
	}

	// Create a scraping task
	task := &domain.ScrapingTask{
		ID:        uuid.New(),
		URLID:     url.ID,
		URL:       url.Url,
		Status:    domain.URLStatusPending,
		Attempt:   1,
		CreatedAt: time.Now(),
	}

	// Send the task to Kafka
	message := domain.NewScrapingTaskMessage(task, uuid.New().String())
	if err := s.producer.SendMessage(ctx, domain.TopicScrapingTasks, message); err != nil {
		return fmt.Errorf("failed to send scraping task to Kafka: %w", err)
	}

	// Update URL status to in_progress
	if err := s.urlRepo.UpdateURLStatus(ctx, url.ID, string(domain.URLStatusInProgress)); err != nil {
		s.logger.WithError(err).WithField("url_id", url.ID).Error("Failed to update URL status")
	}

	// Update last scraped time
	if err := s.urlRepo.UpdateLastScrapedTime(ctx, url.ID, time.Now()); err != nil {
		s.logger.WithError(err).WithField("url_id", url.ID).Error("Failed to update last scraped time")
	}

	// Calculate and update next scrape time
	nextScrape, err := models.CalculateNextScrapeTime(url.Frequency, time.Now())
	if err != nil {
		s.logger.WithError(err).WithField("url_id", url.ID).Error("Failed to calculate next scrape time")
		return err
	}

	if err := s.urlRepo.UpdateNextScrapeTime(ctx, url.ID, nextScrape); err != nil {
		s.logger.WithError(err).WithField("url_id", url.ID).Error("Failed to update next scrape time")
	}

	s.logger.WithFields(logrus.Fields{
		"url_id":      url.ID,
		"url":         url.Url,
		"task_id":     task.ID,
		"next_scrape": nextScrape,
	}).Info("Created scraping task")

	return nil
}
