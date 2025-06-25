# Shared Database Makefile

.PHONY: migrate-up migrate-down migrate-status sqlc-generate sqlc-check

# Database configuration
DATABASE_URL ?= postgres://scraper:scraper@localhost:5432/scraping_db?sslmode=disable
MIGRATIONS_DIR ?= ../../sql/schema
SQLC_CONFIG ?= sqlc.yaml

# Migration commands
migrate-up:
	@echo "Running database migrations..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" up

migrate-down:
	@echo "Rolling back last migration..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" down

migrate-status:
	@echo "Migration status:"
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" status

migrate-reset:
	@echo "Resetting all migrations..."
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DATABASE_URL)" reset

# SQLC commands
sqlc-generate:
	@echo "Generating SQLC code..."
	@sqlc generate -f $(SQLC_CONFIG)

sqlc-check:
	@echo "Checking SQLC configuration..."
	@sqlc compile -f $(SQLC_CONFIG)

sqlc-verify:
	@echo "Verifying SQLC configuration..."
	@sqlc verify -f $(SQLC_CONFIG)

# Database setup
db-setup: migrate-up sqlc-generate
	@echo "Database setup completed"

# Database reset
db-reset: migrate-reset migrate-up sqlc-generate
	@echo "Database reset completed"

# Help
help:
	@echo "Available database commands:"
	@echo "  migrate-up      - Run all pending migrations"
	@echo "  migrate-down    - Rollback last migration"
	@echo "  migrate-status  - Show migration status"
	@echo "  migrate-reset   - Reset all migrations"
	@echo "  sqlc-generate   - Generate Go code from SQL"
	@echo "  sqlc-check      - Check SQLC configuration"
	@echo "  sqlc-verify     - Verify SQLC configuration"
	@echo "  db-setup        - Setup database (migrate + generate)"
	@echo "  db-reset        - Reset database (reset + migrate + generate)" 