package utils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"go_scraping_project/internal/database"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConfig holds test database configuration
type TestDatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// GetTestDatabaseConfig returns test database configuration
func GetTestDatabaseConfig() TestDatabaseConfig {
	return TestDatabaseConfig{
		Host:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		Port:     getEnvOrDefaultInt("TEST_DB_PORT", 5432),
		User:     getEnvOrDefault("TEST_DB_USER", "scraper"),
		Password: getEnvOrDefault("TEST_DB_PASSWORD", "scraper_password"),
		DBName:   getEnvOrDefault("TEST_DB_NAME", "scraper_test"),
		SSLMode:  getEnvOrDefault("TEST_DB_SSLMODE", "disable"),
	}
}

// SetupTestDatabase creates a test database connection
func SetupTestDatabase(t *testing.T) (*sql.DB, func()) {
	config := GetTestDatabaseConfig()

	// Create connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	// Test connection
	err = db.Ping()
	require.NoError(t, err)

	// Return cleanup function
	cleanup := func() {
		db.Close()
	}

	return db, cleanup
}

// CreateTestURL creates a test URL in the database
func CreateTestURL(t *testing.T, db *sql.DB, url, frequency, description string) database.Url {
	querier := database.New(db)

	params := database.CreateURLParams{
		Url:        url,
		Frequency:  frequency,
		Status:     "pending",
		MaxRetries: 3,
		Timeout:    30,
		RateLimit:  1,
	}

	result, err := querier.CreateURL(context.Background(), params)
	require.NoError(t, err)

	return result
}

// CleanupTestData removes test data from database
func CleanupTestData(t *testing.T, db *sql.DB) {
	// Delete all URLs
	_, err := db.Exec("DELETE FROM urls")
	require.NoError(t, err)

	// Reset sequences if needed
	_, err = db.Exec("ALTER SEQUENCE urls_id_seq RESTART WITH 1")
	if err != nil {
		// Ignore error if sequence doesn't exist
		t.Logf("Warning: Could not reset sequence: %v", err)
	}
}

// GenerateTestUUID generates a test UUID
func GenerateTestUUID() uuid.UUID {
	return uuid.New()
}

// GenerateTestURLs creates multiple test URLs
func GenerateTestURLs(t *testing.T, db *sql.DB, count int) []database.Url {
	urls := make([]database.Url, count)

	for i := 0; i < count; i++ {
		url := fmt.Sprintf("https://example%d.com", i+1)
		frequency := "hourly"
		if i%2 == 0 {
			frequency = "daily"
		}
		description := fmt.Sprintf("Test URL %d", i+1)

		urls[i] = CreateTestURL(t, db, url, frequency, description)
	}

	return urls
}

// WaitForCondition waits for a condition to be true with timeout
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}

	return false
}

// AssertURLExists checks if a URL exists in the database
func AssertURLExists(t *testing.T, db *sql.DB, urlID uuid.UUID) {
	querier := database.New(db)

	url, err := querier.GetURLByID(context.Background(), urlID)
	require.NoError(t, err)
	require.NotEmpty(t, url)
	require.Equal(t, urlID, url.ID)
}

// AssertURLNotExists checks if a URL does not exist in the database
func AssertURLNotExists(t *testing.T, db *sql.DB, urlID uuid.UUID) {
	querier := database.New(db)

	_, err := querier.GetURLByID(context.Background(), urlID)
	require.Error(t, err)
}

// GetURLCount returns the total number of URLs in the database
func GetURLCount(t *testing.T, db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM urls").Scan(&count)
	require.NoError(t, err)
	return count
}

// Helper functions for environment variables
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvOrDefaultInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// TestTime utilities
func GetTestTime() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func GetTestTimeString() string {
	return GetTestTime().Format(time.RFC3339)
}

// TestData generators
func GenerateTestURLString() string {
	return fmt.Sprintf("https://example-%s.com", uuid.New().String()[:8])
}

func GenerateTestDescription() string {
	return fmt.Sprintf("Test description %s", uuid.New().String()[:8])
}

// Test context helpers
func CreateTestContext() context.Context {
	return context.Background()
}

func CreateTestContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// Test error helpers
func ExpectError(t *testing.T, err error, expectedErrorType string) {
	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErrorType)
}

func ExpectNoError(t *testing.T, err error) {
	require.NoError(t, err)
}

// Test assertion helpers
func AssertHTTPStatus(t *testing.T, statusCode int, expectedStatus int) {
	require.Equal(t, expectedStatus, statusCode)
}

func AssertResponseContentType(t *testing.T, contentType string, expectedType string) {
	require.Contains(t, contentType, expectedType)
}

// Test cleanup helpers
func CleanupTestFiles(t *testing.T, filePaths ...string) {
	for _, path := range filePaths {
		if _, err := os.Stat(path); err == nil {
			err := os.Remove(path)
			if err != nil {
				t.Logf("Warning: Could not remove test file %s: %v", path, err)
			}
		}
	}
}

// Test logging helpers
func SetupTestLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Only show errors in tests
	return logger
}
