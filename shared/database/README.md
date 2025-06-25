# Shared Database Package

This package provides shared database functionality for all services in the Go Scraping Project.

## Overview

The shared database package includes:
- Database connection management
- Migration tools (Goose)
- SQLC code generation
- Repository interfaces
- Common database operations

## Structure

```
shared/database/
├── connection.go      # Database connection management
├── migrations.go      # Migration utilities
├── repository.go      # Repository interfaces
├── sqlc.yaml         # SQLC configuration
├── Makefile          # Database operations
└── README.md         # This file
```

## Usage

### Connecting to Database

```go
import "go_scraping_project/shared/database"

// Connect to database
db, err := database.Connect()
if err != nil {
    log.Fatal(err)
}
defer db.Close()
```

### Running Migrations

```go
// Run migrations
err := database.RunMigrations(db)
if err != nil {
    log.Fatal(err)
}
```

### Using Repository Interface

```go
// Create repository
repo := database.NewBaseRepository(db)

// Use repository methods
urls, err := repo.GetAllURLs(ctx)
if err != nil {
    log.Fatal(err)
}
```

## Database Operations

### From Top-Level Makefile

```bash
# Setup database (migrate + generate)
make db-setup

# Run migrations
make db-migrate-up

# Rollback last migration
make db-migrate-down

# Show migration status
make db-migrate-status

# Generate SQLC code
make sqlc-generate

# Check SQLC configuration
make sqlc-check
```

### From Database Directory

```bash
cd shared/database

# Setup database
make db-setup

# Run migrations
make migrate-up

# Generate code
make sqlc-generate
```

## Configuration

### Environment Variables

- `DATABASE_URL` - Full PostgreSQL connection string
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: scraper)
- `DB_PASSWORD` - Database password (default: scraper)
- `DB_NAME` - Database name (default: scraping_db)
- `DB_SSLMODE` - SSL mode (default: disable)
- `MIGRATIONS_DIR` - Migrations directory (default: ../../sql/schema)

### Example Connection String

```
postgres://scraper:scraper@localhost:5432/scraping_db?sslmode=disable
```

## SQLC Configuration

The `sqlc.yaml` file configures SQLC code generation:

- **Queries**: `sql/queries/` - SQL query files
- **Schema**: `sql/schema/` - Database schema files
- **Output**: `db/` - Generated Go code
- **Package**: `db` - Go package name

## Migration Files

Migrations are stored in `sql/schema/` and follow the naming convention:
```
001_create_urls_table.sql
002_create_scraping_tasks_table.sql
003_create_scraped_data_table.sql
```

## Generated Code

SQLC generates Go code in the `db/` directory:
- `db.go` - Database connection and queries
- `models.go` - Generated structs
- `querier.go` - Query interface

## Best Practices

1. **Connection Pooling**: The connection is configured with appropriate pool settings
2. **Error Handling**: Always check for errors and handle them appropriately
3. **Context Usage**: Use context for cancellation and timeouts
4. **Transaction Management**: Use transactions for multi-step operations
5. **Migration Safety**: Always backup before running migrations in production

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check if PostgreSQL is running
   - Verify connection string
   - Check firewall settings

2. **Migration Failures**
   - Check migration file syntax
   - Verify database permissions
   - Check for conflicting migrations

3. **SQLC Generation Errors**
   - Verify SQL syntax
   - Check sqlc.yaml configuration
   - Ensure schema files are up to date

### Debugging

```bash
# Check database connection
psql "$DATABASE_URL" -c "SELECT 1;"

# Check migration status
make db-migrate-status

# Validate SQLC configuration
make sqlc-check
``` 