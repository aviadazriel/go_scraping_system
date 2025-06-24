package database

import (
	"context"
	"database/sql"
	"fmt"

	"go_scraping_project/internal/config"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// PostgresDB represents a PostgreSQL database connection for sqlc
type PostgresDB struct {
	DB     *sql.DB
	config *config.Config
	logger *logrus.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection for sqlc
func NewPostgresDB(cfg *config.Config, log *logrus.Logger) (*PostgresDB, error) {
	dsn := cfg.GetDatabaseURL()

	// Open database connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	postgresDB := &PostgresDB{
		DB:     db,
		config: cfg,
		logger: log,
	}

	log.Info("Database connection established successfully")
	return postgresDB, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}

// HealthCheck performs a database health check
func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	return p.DB.PingContext(ctx)
}

// GetStats returns database connection statistics
func (p *PostgresDB) GetStats() map[string]interface{} {
	stats := p.DB.Stats()

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// GetDB returns the underlying *sql.DB for sqlc
func (p *PostgresDB) GetDB() *sql.DB {
	return p.DB
}
