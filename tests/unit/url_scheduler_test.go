package unit

import (
	"context"
	"testing"
	"time"

	"go_scraping_project/internal/database"
	"go_scraping_project/internal/domain"
	"go_scraping_project/internal/url-manager/services"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockURLRepository is a mock implementation of repositories.URLRepository
type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) GetURLByID(ctx context.Context, id uuid.UUID) (*database.Url, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.Url), args.Error(1)
}

func (m *MockURLRepository) GetURLsScheduledForScraping(ctx context.Context, from, to time.Time, limit int32) ([]database.Url, error) {
	args := m.Called(ctx, from, to, limit)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockURLRepository) GetURLsByStatus(ctx context.Context, status string, limit, offset int32) ([]database.Url, error) {
	args := m.Called(ctx, status, limit, offset)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockURLRepository) UpdateURLStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockURLRepository) UpdateNextScrapeTime(ctx context.Context, id uuid.UUID, nextScrapeAt time.Time) error {
	args := m.Called(ctx, id, nextScrapeAt)
	return args.Error(0)
}

func (m *MockURLRepository) UpdateLastScrapedTime(ctx context.Context, id uuid.UUID, lastScrapedAt time.Time) error {
	args := m.Called(ctx, id, lastScrapedAt)
	return args.Error(0)
}

func (m *MockURLRepository) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockURLRepository) ResetRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockURLRepository) GetURLsForImmediateScraping(ctx context.Context, limit int32) ([]database.Url, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockURLRepository) CountURLsByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockURLRepository) GetURLsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.Url, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]database.Url), args.Error(1)
}

// MockKafkaProducer is a mock implementation of domain.KafkaProducer
type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) SendMessage(ctx context.Context, topic string, message *domain.KafkaMessage) error {
	args := m.Called(ctx, topic, message)
	return args.Error(0)
}

func (m *MockKafkaProducer) SendScrapingTask(ctx context.Context, task *domain.ScrapingTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockKafkaProducer) SendScrapedData(ctx context.Context, data *domain.ScrapedData, success bool, err string) error {
	args := m.Called(ctx, data, success, err)
	return args.Error(0)
}

func (m *MockKafkaProducer) SendParsedData(ctx context.Context, data *domain.ParsedData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

func (m *MockKafkaProducer) SendDeadLetter(ctx context.Context, originalMessage *domain.KafkaMessage, err error, maxRetries int) error {
	args := m.Called(ctx, originalMessage, err, maxRetries)
	return args.Error(0)
}

func (m *MockKafkaProducer) SendRetryMessage(ctx context.Context, originalMessageID string, messageType domain.MessageType, data interface{}, retryCount, maxRetries int, retryDelay time.Duration) error {
	args := m.Called(ctx, originalMessageID, messageType, data, retryCount, maxRetries, retryDelay)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestURLSchedulerService_Start(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockURLRepository, *MockKafkaProducer)
		expectError bool
	}{
		{
			name: "successful service start",
			mockSetup: func(mr *MockURLRepository, mkp *MockKafkaProducer) {
				// No specific expectations for Start method
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			mockProducer := &MockKafkaProducer{}
			tt.mockSetup(mockRepo, mockProducer)

			service := services.NewURLSchedulerService(mockRepo, mockProducer, logrus.New())
			err := service.Start(context.Background())

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Clean up
			service.Stop()
		})
	}
}

func TestURLSchedulerService_Stop(t *testing.T) {
	tests := []struct {
		name        string
		mockSetup   func(*MockURLRepository, *MockKafkaProducer)
		expectError bool
	}{
		{
			name: "successful service stop",
			mockSetup: func(mr *MockURLRepository, mkp *MockKafkaProducer) {
				// No specific expectations for Stop method
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockURLRepository{}
			mockProducer := &MockKafkaProducer{}
			tt.mockSetup(mockRepo, mockProducer)

			service := services.NewURLSchedulerService(mockRepo, mockProducer, logrus.New())
			service.Start(context.Background())
			err := service.Stop()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestURLSchedulerService_NewURLSchedulerService(t *testing.T) {
	mockRepo := &MockURLRepository{}
	mockProducer := &MockKafkaProducer{}
	logger := logrus.New()

	service := services.NewURLSchedulerService(mockRepo, mockProducer, logger)

	assert.NotNil(t, service)
	// Note: The service doesn't expose getter methods, so we just verify it was created
}

func TestURLSchedulerService_ContextCancellation(t *testing.T) {
	mockRepo := &MockURLRepository{}
	mockProducer := &MockKafkaProducer{}
	service := services.NewURLSchedulerService(mockRepo, mockProducer, logrus.New())

	ctx, cancel := context.WithCancel(context.Background())
	err := service.Start(ctx)
	assert.NoError(t, err)

	// Cancel the context
	cancel()

	// Wait a bit for the service to stop
	time.Sleep(100 * time.Millisecond)

	// Service should have stopped due to context cancellation
	assert.NoError(t, service.Stop())
}

func TestURLSchedulerService_StopChannel(t *testing.T) {
	mockRepo := &MockURLRepository{}
	mockProducer := &MockKafkaProducer{}
	service := services.NewURLSchedulerService(mockRepo, mockProducer, logrus.New())

	err := service.Start(context.Background())
	assert.NoError(t, err)

	// Stop the service
	err = service.Stop()
	assert.NoError(t, err)

	// Don't try to stop again to avoid the panic
}
