package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_scraping_project/internal/api-gateway/models"
	"go_scraping_project/internal/api-gateway/types"
	"go_scraping_project/internal/database"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuerier is a mock implementation of database.Querier for testing
type MockQuerier struct {
	mock.Mock
}

func (m *MockQuerier) GetURLByID(ctx context.Context, id uuid.UUID) (database.Url, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(database.Url), args.Error(1)
}

func (m *MockQuerier) GetURLsScheduledForScraping(ctx context.Context, arg database.GetURLsScheduledForScrapingParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockQuerier) GetURLsByStatus(ctx context.Context, arg database.GetURLsByStatusParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockQuerier) UpdateURLStatus(ctx context.Context, arg database.UpdateURLStatusParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdateNextScrapeTime(ctx context.Context, arg database.UpdateNextScrapeTimeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) UpdateLastScrapedTime(ctx context.Context, arg database.UpdateLastScrapedTimeParams) error {
	args := m.Called(ctx, arg)
	return args.Error(0)
}

func (m *MockQuerier) IncrementRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) ResetRetryCount(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockQuerier) GetURLsForImmediateScraping(ctx context.Context, arg database.GetURLsForImmediateScrapingParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockQuerier) CountURLsByStatus(ctx context.Context, status string) (int64, error) {
	args := m.Called(ctx, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockQuerier) GetURLsByIDs(ctx context.Context, ids []uuid.UUID) ([]database.Url, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockQuerier) CreateURL(ctx context.Context, arg database.CreateURLParams) (database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(database.Url), args.Error(1)
}

func (m *MockQuerier) ListURLs(ctx context.Context, arg database.ListURLsParams) ([]database.Url, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]database.Url), args.Error(1)
}

func (m *MockQuerier) CountURLs(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// setupTestHandler creates a test handler with a real database connection
func setupTestHandler() *types.URLHandler {
	logger := logrus.New()

	// Create handler with real database (for integration testing)
	handler := types.NewURLHandler(logger, &database.Queries{})

	return handler
}

func TestURLHandler_CreateURL_Validation_Integration(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    models.CreateURLRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "invalid request body - empty URL",
			requestBody: models.CreateURLRequest{
				URL: "", // Invalid empty URL
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - empty frequency",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "", // Invalid empty frequency
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - invalid URL format",
			requestBody: models.CreateURLRequest{
				URL:       "not-a-url",
				Frequency: "1h",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - invalid frequency format",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "invalid-frequency",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - negative timeout",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "1h",
				Timeout:   -1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - timeout too high",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "1h",
				Timeout:   301, // Max is 300
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - negative rate limit",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "1h",
				RateLimit: -1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - rate limit too high",
			requestBody: models.CreateURLRequest{
				URL:       "https://example.com",
				Frequency: "1h",
				RateLimit: 1001, // Max is 1000
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - negative max retries",
			requestBody: models.CreateURLRequest{
				URL:        "https://example.com",
				Frequency:  "1h",
				MaxRetries: -1,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "invalid request body - max retries too high",
			requestBody: models.CreateURLRequest{
				URL:        "https://example.com",
				Frequency:  "1h",
				MaxRetries: 11, // Max is 10
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, _ := json.Marshal(tt.requestBody)

			// Create HTTP request
			req := httptest.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler directly
			handler.CreateURL(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				// Should return error message
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}

func TestURLHandler_GetURL_Validation_Integration(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		urlID          string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "invalid UUID - too short",
			urlID:          "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid UUID - wrong format",
			urlID:          "12345678-1234-1234-1234-123456789012",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid UUID - empty",
			urlID:          "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest("GET", "/api/v1/urls/"+tt.urlID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler directly
			handler.GetURL(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				// Should return error message
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}

func TestURLHandler_UpdateURL_Validation_Integration(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		urlID          string
		requestBody    models.UpdateURLRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name:  "invalid UUID",
			urlID: "invalid-uuid",
			requestBody: models.UpdateURLRequest{
				Frequency: "1d",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:  "invalid frequency format",
			urlID: "12345678-1234-1234-1234-123456789012",
			requestBody: models.UpdateURLRequest{
				Frequency: "invalid-frequency",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request body
			body, _ := json.Marshal(tt.requestBody)

			// Create HTTP request
			req := httptest.NewRequest("PUT", "/api/v1/urls/"+tt.urlID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler directly
			handler.UpdateURL(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				// Should return error message
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}

func TestURLHandler_DeleteURL_Validation_Integration(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		urlID          string
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "invalid UUID",
			urlID:          "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "invalid UUID - empty",
			urlID:          "",
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP request
			req := httptest.NewRequest("DELETE", "/api/v1/urls/"+tt.urlID, nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Call handler directly
			handler.DeleteURL(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError {
				// Should return error message
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}
