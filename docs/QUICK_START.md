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

# Install dependencies
go mod tidy
```

## Step 2: Database Setup

```bash
# Create database user and database
psql -h localhost -U aazriel -d postgres -c "CREATE USER scraper WITH PASSWORD 'scraper_password';"
psql -h localhost -U aazriel -d postgres -c "CREATE DATABASE scraper OWNER scraper;"

# Run database migrations
DATABASE_URL="postgres://scraper:scraper_password@localhost:5432/scraper?sslmode=disable" make migrate-up
```

## Step 3: Deploy Services

```bash
# Start all services
docker-compose -f docker-compose.local.yml up -d

# Wait for services to start (about 30 seconds)
sleep 30

# Check service status
docker-compose -f docker-compose.local.yml ps
```

## Step 4: Verify Deployment

```bash
# Check API Gateway health
curl http://localhost:8082/health

# Check URL Manager logs
docker-compose -f docker-compose.local.yml logs url-manager
```

## Step 5: Test the API

```bash
# Create a test URL
curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "frequency": "hourly",
    "description": "Test website"
  }'

# List all URLs
curl http://localhost:8082/api/v1/urls
```

## Step 6: Monitor the System

```bash
# View all service logs
docker-compose -f docker-compose.local.yml logs -f

# Open Kafka UI in your browser
open http://localhost:8080
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

3. **Check logs:**
   ```bash
   docker-compose -f docker-compose.local.yml logs
   ```

### If you get build errors:

```bash
# Update dependencies
go mod tidy

# Rebuild
docker-compose -f docker-compose.local.yml build --no-cache
```

## Next Steps

- Read the full [README.md](../README.md) for detailed documentation
- Check out the [API Reference](../README.md#api-reference)
- Explore the [Example Workflows](../README.md#example-workflows)

## Ports Used

| Service | Port | Description |
|---------|------|-------------|
| API Gateway | 8082 | REST API |
| Kafka UI | 8080 | Kafka monitoring |
| Kafka | 9092 | Kafka broker |
| URL Manager | 8081 | Background service |

## Useful Commands

```bash
# Stop all services
docker-compose -f docker-compose.local.yml down

# Restart a specific service
docker-compose -f docker-compose.local.yml restart url-manager

# View logs for a specific service
docker-compose -f docker-compose.local.yml logs -f api-gateway

# Check database
psql -h localhost -U scraper -d scraper -c "SELECT * FROM urls;"
``` 