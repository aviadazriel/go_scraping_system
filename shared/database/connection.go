package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database
func Connect() (*sql.DB, error) {
	// Get database URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Fallback to individual environment variables
		host := getEnvOrDefault("DB_HOST", "localhost")
		port := getEnvOrDefault("DB_PORT", "5432")
		user := getEnvOrDefault("DB_USER", "scraper")
		password := getEnvOrDefault("DB_PASSWORD", "scraper")
		dbName := getEnvOrDefault("DB_NAME", "scraping_db")
		sslMode := getEnvOrDefault("DB_SSLMODE", "disable")

		databaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbName, sslMode)
	}

	// Open database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// ConnectWithConfig establishes a connection using configuration
func ConnectWithConfig(cfg interface{}) (*sql.DB, error) {
	// Type assertion to get config methods
	config, ok := cfg.(interface {
		GetString(key string) string
		GetInt(key string) int
	})
	if !ok {
		return nil, fmt.Errorf("config does not implement required methods")
	}

	// Build database URL from config
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	dbName := config.GetString("database.database")
	sslMode := config.GetString("database.ssl_mode")

	// Fallback to environment variables if config values are empty
	if host == "" {
		host = getEnvOrDefault("DB_HOST", "localhost")
	}
	if port == 0 {
		portStr := getEnvOrDefault("DB_PORT", "5432")
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		} else {
			port = 5432
		}
	}
	if user == "" {
		user = getEnvOrDefault("DB_USER", "scraper")
	}
	if password == "" {
		password = getEnvOrDefault("DB_PASSWORD", "scraper")
	}
	if dbName == "" {
		dbName = getEnvOrDefault("DB_NAME", "scraping_db")
	}
	if sslMode == "" {
		sslMode = getEnvOrDefault("DB_SSLMODE", "disable")
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

	// Open database connection
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool from config
	maxOpenConns := config.GetInt("database.max_open_conns")
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	db.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := config.GetInt("database.max_idle_conns")
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}
	db.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime := config.GetString("database.conn_max_lifetime")
	if connMaxLifetime == "" {
		connMaxLifetime = "5m"
	}
	if duration, err := time.ParseDuration(connMaxLifetime); err == nil {
		db.SetConnMaxLifetime(duration)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// getEnvOrDefault returns the environment variable value or a default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Close closes the database connection
func Close(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}
