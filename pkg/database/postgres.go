package database

import (
	"context"
	"fmt"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/internal/domain"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB represents a PostgreSQL database connection
type PostgresDB struct {
	DB     *gorm.DB
	config *config.Config
	logger *logrus.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config, log *logrus.Logger) (*PostgresDB, error) {
	dsn := cfg.GetDatabaseURL()
	
	// Configure GORM logger
	gormLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	postgresDB := &PostgresDB{
		DB:     db,
		config: cfg,
		logger: log,
	}

	// Run migrations
	if err := postgresDB.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return postgresDB, nil
}

// Migrate runs database migrations
func (p *PostgresDB) Migrate() error {
	p.logger.Info("Running database migrations...")

	// Auto-migrate all models
	err := p.DB.AutoMigrate(
		&domain.URL{},
		&domain.ScrapingTask{},
		&domain.ScrapedData{},
		&domain.ParsedData{},
		&domain.DeadLetterMessage{},
		&domain.ScrapingMetrics{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate: %w", err)
	}

	p.logger.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// HealthCheck performs a database health check
func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.PingContext(ctx)
}

// GetStats returns database connection statistics
func (p *PostgresDB) GetStats() map[string]interface{} {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return map[string]interface{}{
			"error": "failed to get underlying sql.DB",
		}
	}

	return map[string]interface{}{
		"max_open_connections": sqlDB.Stats().MaxOpenConnections,
		"open_connections":     sqlDB.Stats().OpenConnections,
		"in_use":              sqlDB.Stats().InUse,
		"idle":                sqlDB.Stats().Idle,
		"wait_count":          sqlDB.Stats().WaitCount,
		"wait_duration":       sqlDB.Stats().WaitDuration,
		"max_idle_closed":     sqlDB.Stats().MaxIdleClosed,
		"max_lifetime_closed": sqlDB.Stats().MaxLifetimeClosed,
	}
} 