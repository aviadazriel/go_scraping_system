package unit

import (
	"context"
	"testing"
	"time"

	"go_scraping_project/internal/database"
	"go_scraping_project/internal/url-manager/repositories"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockQuerier is a mock implementation of database.Querier
type mockQuerier struct {
	mock.Mock
}

func (m *mockQuerier) GetURLByID(ctx context.Context, id uuid.UUID) (database.Url, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Url), args.Error(1)
}

func (m *mockQuerier) GetURLsScheduledForScraping(ctx context.Context, arg database.GetURLsScheduledForScrapingParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *mockQuerier) GetURLsByStatus(ctx context.Context, arg database.GetURLsByStatusParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *mockQuerier) UpdateURLStatus(ctx context.Context, arg database.UpdateURLStatusParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) UpdateNextScrapeTime(ctx context.Context, arg database.UpdateNextScrapeTimeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) UpdateLastScrapedTime(ctx context.Context, arg database.UpdateLastScrapedTimeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *mockQuerier) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockQuerier) ResetRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockQuerier) GetURLsForImmediateScraping(ctx context.Context, arg database.GetURLsForImmediateScrapingParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *mockQuerier) CountURLsByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockQuerier) GetURLsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.Url, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]database.Url), args.Error(1)
}

func TestURLRepository_GetURLByID(t *testing.T) {
	// Setup
	mockDB := &mockQuerier{}
	logger := logrus.New()
	repo := repositories.NewURLRepository(mockDB, logger)

	ctx := context.Background()
	urlID := uuid.New()
	expectedURL := database.Url{
		ID:        urlID,
		Url:       "https://example.com",
		Frequency: "daily",
		Status:    "active",
	}

	// Expectations
	mockDB.On("GetURLByID", ctx, urlID).Return(expectedURL, nil)

	// Execute
	result, err := repo.GetURLByID(ctx, urlID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL.ID, result.ID)
	assert.Equal(t, expectedURL.Url, result.Url)
	assert.Equal(t, expectedURL.Frequency, result.Frequency)
	assert.Equal(t, expectedURL.Status, result.Status)

	mockDB.AssertExpectations(t)
}

func TestURLRepository_GetURLsScheduledForScraping(t *testing.T) {
	// Setup
	mockDB := &mockQuerier{}
	logger := logrus.New()
	repo := repositories.NewURLRepository(mockDB, logger)

	ctx := context.Background()
	from := time.Now().UTC()
	to := from.Add(time.Hour)
	limit := int32(10)

	expectedURLs := []database.Url{
		{
			ID:        uuid.New(),
			Url:       "https://example1.com",
			Frequency: "daily",
			Status:    "active",
		},
		{
			ID:        uuid.New(),
			Url:       "https://example2.com",
			Frequency: "hourly",
			Status:    "active",
		},
	}

	// Expectations
	mockDB.On("GetURLsScheduledForScraping", ctx, mock.AnythingOfType("database.GetURLsScheduledForScrapingParams")).Return(expectedURLs, nil)

	// Execute
	result, err := repo.GetURLsScheduledForScraping(ctx, from, to, limit)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedURLs[0].ID, result[0].ID)
	assert.Equal(t, expectedURLs[1].ID, result[1].ID)

	mockDB.AssertExpectations(t)
}

func TestURLRepository_UpdateURLStatus(t *testing.T) {
	// Setup
	mockDB := &mockQuerier{}
	logger := logrus.New()
	repo := repositories.NewURLRepository(mockDB, logger)

	ctx := context.Background()
	urlID := uuid.New()
	newStatus := "processing"

	// Expectations
	mockDB.On("UpdateURLStatus", ctx, mock.AnythingOfType("database.UpdateURLStatusParams")).Return(nil)

	// Execute
	err := repo.UpdateURLStatus(ctx, urlID, newStatus)

	// Assert
	assert.NoError(t, err)

	mockDB.AssertExpectations(t)
}

func TestURLRepository_GetURLsForImmediateScraping(t *testing.T) {
	// Setup
	mockDB := &mockQuerier{}
	logger := logrus.New()
	repo := repositories.NewURLRepository(mockDB, logger)

	ctx := context.Background()
	limit := int32(5)

	expectedURLs := []database.Url{
		{
			ID:        uuid.New(),
			Url:       "https://example1.com",
			Frequency: "immediate",
			Status:    "pending",
		},
	}

	// Expectations
	mockDB.On("GetURLsForImmediateScraping", ctx, mock.AnythingOfType("database.GetURLsForImmediateScrapingParams")).Return(expectedURLs, nil)

	// Execute
	result, err := repo.GetURLsForImmediateScraping(ctx, limit)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, expectedURLs[0].ID, result[0].ID)

	mockDB.AssertExpectations(t)
}

func TestURLRepository_CountURLsByStatus(t *testing.T) {
	// Setup
	mockDB := &mockQuerier{}
	logger := logrus.New()
	repo := repositories.NewURLRepository(mockDB, logger)

	ctx := context.Background()
	status := "active"
	expectedCount := int64(42)

	// Expectations
	mockDB.On("CountURLsByStatus", ctx, status).Return(expectedCount, nil)

	// Execute
	result, err := repo.CountURLsByStatus(ctx, status)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)

	mockDB.AssertExpectations(t)
}
