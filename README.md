# Go Scraping Project - Multi-Project Monorepo

A microservices-based web scraping platform built with Go, following clean architecture principles and using Kafka for event-driven communication.

## Project Structure

This project follows a **multi-project monorepo** structure where each service has its own `go.mod`, Dockerfile, and Makefile while sharing common code through a shared package.

```
go_scraping_project/
├── shared/                    # Shared packages used across all services
│   ├── utils/                # Common utilities (time, validation, etc.)
│   ├── models/               # Shared domain models
│   ├── config/               # Shared configuration structures
│   ├── database/             # Shared database functionality
│   │   ├── connection.go     # Database connection management
│   │   ├── migrations.go     # Migration utilities
│   │   ├── repository.go     # Repository interfaces
│   │   ├── sqlc.yaml         # SQLC configuration
│   │   └── Makefile          # Database operations
│   └── go.mod                # Shared dependencies
├── services/                 # Individual microservices
│   ├── api-gateway/          # API Gateway service
│   │   ├── handlers/         # HTTP handlers
│   │   ├── main.go           # Service entry point
│   │   ├── go.mod            # Service-specific dependencies
│   │   ├── Dockerfile        # Service-specific Docker image
│   │   └── Makefile          # Service-specific build commands
│   ├── url-manager/          # URL Manager service
│   │   ├── handlers/         # HTTP handlers
│   │   ├── services/         # Business logic
│   │   ├── repositories/     # Data access layer
│   │   ├── main.go           # Service entry point
│   │   ├── go.mod            # Service-specific dependencies
│   │   ├── Dockerfile        # Service-specific Docker image
│   │   └── Makefile          # Service-specific build commands
│   ├── scraper/              # Web Scraper service (planned)
│   ├── parser/               # Content Parser service (planned)
│   └── storage/              # Data Storage service (planned)
├── sql/                      # Database schema and migrations
│   ├── schema/               # Migration files
│   └── queries/              # SQLC query files
├── monitoring/               # Monitoring configuration
├── docker-compose.yml        # Multi-service orchestration
├── Makefile                  # Top-level build commands
└── README.md                 # This file
```

## Architecture

### Services

1. **API Gateway** (`:8080`) - Central entry point for external clients
2. **URL Manager** (`:8081`) - Manages URLs to be scraped and scheduling
3. **Scraper** (planned) - Performs actual web scraping
4. **Parser** (planned) - Parses and structures scraped content
5. **Storage** (planned) - Stores and retrieves parsed data

### Infrastructure

- **PostgreSQL** (`:5432`) - Primary database
- **Kafka** (`:9092`) - Message broker for event-driven communication
- **Zookeeper** (`:2181`) - Required for Kafka
- **Kafka UI** (`:8080`) - Web interface for Kafka management
- **Prometheus** (`:9090`) - Metrics collection
- **Grafana** (`:3000`) - Metrics visualization

## Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Make
- PostgreSQL (for local development)

### Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go_scraping_project
   ```

2. **Setup database**
   ```bash
   # Start PostgreSQL
   docker-compose up -d postgres
   
   # Setup database (migrations + code generation)
   make db-setup
   ```

3. **Setup all services**
   ```bash
   make dev-setup-all
   ```

4. **Start infrastructure**
   ```bash
   docker-compose up -d zookeeper kafka
   ```

5. **Run services individually**
   ```bash
   # Run API Gateway
   make api-gateway
   
   # Run URL Manager (in another terminal)
   make url-manager
   ```

### Docker Deployment

1. **Build and run all services**
   ```bash
   make docker-build-all
   docker-compose up -d
   ```

2. **View logs**
   ```bash
   docker-compose logs -f
   ```

3. **Stop all services**
   ```bash
   docker-compose down
   ```

## Service Development

### Working with Individual Services

Each service can be developed independently:

```bash
# Navigate to a service
cd services/api-gateway

# Install dependencies
make deps

# Run tests
make test

# Build the service
make build

# Run the service
make run

# Format code
make fmt

# Lint code
make lint
```

### Adding a New Service

1. Create a new directory in `services/`
2. Copy the structure from an existing service
3. Update the service-specific `go.mod` with required dependencies
4. Add the service to the top-level `Makefile`
5. Update `docker-compose.yml` if needed

### Shared Code

Common utilities and models are in the `shared/` directory:

- `shared/utils/` - Time utilities, validation, etc.
- `shared/models/` - Domain models used across services
- `shared/config/` - Configuration structures
- `shared/database/` - Database connection, migrations, and repository interfaces

## Database Operations

### Database Setup

```bash
# Setup database (migrations + code generation)
make db-setup

# Run migrations only
make db-migrate-up

# Rollback last migration
make db-migrate-down

# Show migration status
make db-migrate-status

# Generate SQLC code only
make sqlc-generate

# Check SQLC configuration
make sqlc-check
```

### Database Configuration

The database is shared across all services and configured through environment variables:

- `DATABASE_URL` - Full PostgreSQL connection string
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` - Individual connection parameters
- `MIGRATIONS_DIR` - Migrations directory (default: sql/schema)

### Using Shared Database

```go
import "go_scraping_project/shared/database"

// Connect to database
db, err := database.Connect()
if err != nil {
    log.Fatal(err)
}
defer db.Close()

// Run migrations
err = database.RunMigrations(db)
if err != nil {
    log.Fatal(err)
}
```

## API Endpoints

### API Gateway (`:8080`)

- `GET /health` - Health check
- `GET /urls` - List all URLs
- `POST /urls` - Create a new URL
- `GET /urls/{id}` - Get URL by ID
- `PUT /urls/{id}` - Update URL
- `DELETE /urls/{id}` - Delete URL

### URL Manager (`:8081`)

- `GET /health` - Health check
- `GET /urls` - List all URLs
- `POST /urls` - Create a new URL
- `GET /urls/{id}` - Get URL by ID
- `PUT /urls/{id}` - Update URL
- `DELETE /urls/{id}` - Delete URL
- `POST /urls/{id}/schedule` - Schedule URL for scraping

## Configuration

### Environment Variables

Each service can be configured using environment variables:

- `DATABASE_URL` - PostgreSQL connection string
- `KAFKA_BROKERS` - Comma-separated list of Kafka brokers
- `LOG_LEVEL` - Logging level (debug, info, warn, error)

### Database Configuration

The database schema is automatically applied when the PostgreSQL container starts. Schema files are located in `sql/schema/`.

## Testing

### Running Tests

```bash
# Test all services
make test-all

# Test specific service
cd services/api-gateway && make test

# Test with coverage
cd services/api-gateway && make test-coverage
```

### Test Structure

Each service has its own test files:
- Unit tests in the same package as the code
- Integration tests in `*_test.go` files
- Mock implementations for external dependencies

## Monitoring

### Metrics

Services expose Prometheus metrics on `/metrics` endpoints.

### Logging

All services use structured JSON logging with correlation IDs for request tracing.

### Health Checks

Each service provides a `/health` endpoint for monitoring.

## Development Workflow

1. **Feature Development**
   - Create a feature branch
   - Develop in the specific service directory
   - Write tests for new functionality
   - Update shared code if needed

2. **Database Changes**
   - Create migration files in `sql/schema/`
   - Update SQLC queries in `sql/queries/`
   - Run `make db-setup` to apply changes
   - Update repository implementations if needed

3. **Testing**
   - Run unit tests: `make test`
   - Run integration tests: `make test-integration`
   - Check code quality: `make lint`

4. **Building**
   - Build individual service: `make build`
   - Build all services: `make build-all`

5. **Deployment**
   - Build Docker images: `make docker-build-all`
   - Deploy with Docker Compose: `docker-compose up -d`

## Contributing

1. Follow Go best practices and the project's coding standards
2. Write tests for all new functionality
3. Update documentation as needed
4. Use conventional commit messages
5. Ensure all tests pass before submitting a PR

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 8080, 8081, 5432, 9092 are available
2. **Database connection**: Wait for PostgreSQL to fully start before running services
3. **Kafka connection**: Ensure Zookeeper is running before starting Kafka
4. **Migration failures**: Check migration file syntax and database permissions

### Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f api-gateway
docker-compose logs -f url-manager
```

### Database Issues

```bash
# Check database connection
psql "$DATABASE_URL" -c "SELECT 1;"

# Check migration status
make db-migrate-status

# Reset database (careful!)
make db-reset
```

## License

[Add your license information here] 