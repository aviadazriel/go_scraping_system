# Quick Start Guide

This guide will get you up and running with the Go Scraping Project in under 10 minutes.

## Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (running locally on port 5432)

## Step 1: Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd go_scraping_project

# Install dependencies for all services
make deps
```

## Step 2: Database Setup

```bash
# Start PostgreSQL and Kafka using Docker Compose
docker-compose up -d postgres kafka zookeeper

# Wait for services to be ready (about 30 seconds)
sleep 30

# Run database migrations
make migrate-up
```

## Step 3: Create Required Kafka Topics

Before running the services, ensure the required Kafka topics exist:

```bash
docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic scraping-tasks --partitions 1 --replication-factor 1
# Add more topics as needed (e.g. scraping-results, url-updates)
```

## Step 4: Configuration

The project uses a shared configuration system with inheritance:

- `configs/shared.yaml` - Base configuration for all services
- `configs/api-gateway.yaml` - API Gateway specific settings
- `configs/url-manager.yaml` - URL Manager specific settings

You can override settings using environment variables:
```bash
export SCRAPER_DATABASE_HOST=localhost
export SCRAPER_KAFKA_BROKERS=localhost:9092
export SCRAPER_LOG_LEVEL=info
```

## Step 5: Build and Deploy Services

### Option A: Deploy with Docker Compose (Recommended)

```bash
# Build and start all services
docker-compose -f docker-compose.local.yml up -d --build

# Wait for services to start (about 30 seconds)
sleep 30

# Check service status
docker-compose -f docker-compose.local.yml ps
```

### Option B: Run Services Locally

```bash
# Build all services
make build

# Run API Gateway
cd services/api-gateway
./api-gateway

# In another terminal, run URL Manager
cd services/url-manager
./url-manager
```

## Step 6: Verify Deployment

```bash
# Check API Gateway health
curl http://localhost:8082/health

# Check URL Manager logs
docker-compose -f docker-compose.local.yml logs url-manager

# Or if running locally:
# Check URL Manager is running
ps aux | grep url-manager
```

## Step 7: Test the API

```bash
# Create a test URL
curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "frequency": "1h",
    "description": "Test website"
  }'

# List all URLs
curl http://localhost:8082/api/v1/urls

# Check URL status
curl http://localhost:8082/api/v1/urls/{url_id}
```

## Step 8: Monitor the System

```bash
# View all service logs
docker-compose -f docker-compose.local.yml logs -f

# Open Kafka UI in your browser
open http://localhost:8080

# Check database
psql -h localhost -U scraper -d scraping_db -c "SELECT * FROM urls;"
```

## Step 9: Run Tests

```bash
# Run all tests
make test

# Run specific service tests
cd services/url-manager && go test -v
cd services/api-gateway && go test -v
```

## Troubleshooting

### If services fail to start:

1. **Check PostgreSQL:**
   ```bash
   pg_isready -h localhost -p 5432
   ```

2. **Check Docker:**
   ```bash
   docker --version
   docker-compose --version
   ```

3. **Check configuration:**
   ```bash
   # Verify config files exist
   ls -la configs/
   
   # Test configuration loading
   cd services/url-manager && go test -v -run TestConfigLoading
   ```

4. **Check logs:**
   ```bash
   docker-compose -f docker-compose.local.yml logs
   ```

### If you get build errors:

```bash
# Update dependencies
make deps

# Clean and rebuild
make clean
make build

# Or for Docker
docker-compose -f docker-compose.local.yml build --no-cache
```

### If database connection fails:

```bash
# Check database URL
echo $DATABASE_URL

# Test database connection
psql $DATABASE_URL -c "SELECT 1;"
```

## Project Structure

```
go_scraping_project/
├── shared/                 # Shared packages (config, database, kafka)
├── services/              # Microservices
│   ├── api-gateway/       # REST API service
│   ├── url-manager/       # URL scheduling service
│   ├── scraper/           # Web scraping service
│   ├── parser/            # Data parsing service
│   └── storage/           # Data storage service
├── configs/               # Configuration files
├── sql/                   # Database schema and migrations
├── docs/                  # Documentation
└── tests/                 # Integration tests
```

## Ports Used

| Service | Port | Description |
|---------|------|-------------|
| API Gateway | 8082 | REST API |
| Kafka UI | 8080 | Kafka monitoring |
| Kafka | 9092 | Kafka broker |
| PostgreSQL | 5432 | Database |
| URL Manager | 8081 | Background service (if running locally) |

## Useful Commands

```bash
# Stop all services
docker-compose -f docker-compose.local.yml down

# Restart a specific service
docker-compose -f docker-compose.local.yml restart url-manager

# View logs for a specific service
docker-compose -f docker-compose.local.yml logs -f api-gateway

# Check database
psql -h localhost -U scraper -d scraping_db -c "SELECT * FROM urls;"

# Run database migrations
make migrate-up

# Generate SQLC code
make sqlc-generate

# Format code
make fmt

# Lint code
make lint
```

## Next Steps

- Read the full [README.md](../README.md) for detailed documentation
- Check out the [API Reference](../README.md#api-reference)
- Explore the [Example Workflows](../README.md#example-workflows)
- Review the [Architecture Documentation](../docs/ARCHITECTURE.md) 