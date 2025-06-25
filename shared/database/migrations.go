package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pressly/goose/v3"
)

// RunMigrations runs database migrations using goose
func RunMigrations(db *sql.DB) error {
	// Set goose configuration
	goose.SetBaseFS(nil) // Use local filesystem

	// Get migrations directory from environment or use default
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "sql/schema"
	}

	// Run migrations
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// RollbackMigrations rolls back the last migration
func RollbackMigrations(db *sql.DB) error {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "sql/schema"
	}

	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func GetMigrationStatus(db *sql.DB) error {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "sql/schema"
	}

	if err := goose.Status(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to get migration status: %w", err)
	}

	return nil
}
