# Go Scraping Project

A microservices-based web scraping platform built with Go, featuring event-driven architecture using Kafka, PostgreSQL for data persistence, and clean architecture principles.

## 🏗️ Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │  URL Manager    │    │   Scraper       │
│   (Port 8082)   │◄──►│  (Background)   │───►│   (Future)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │     Kafka       │    │   Parser        │
│   (Local)       │    │   (Port 9092)   │    │   (Future)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Services

- **API Gateway** (`:8082`): REST API for managing URLs and viewing data
- **URL Manager** (Background): Schedules and triggers scraping tasks
- **Kafka** (`:9092`): Message broker for service communication
- **PostgreSQL** (Local): Data persistence
- **Kafka UI** (`:8080`): Web interface for monitoring Kafka

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (running locally on port 5432)

### 1. Setup Database

```bash
# Create database user and database
psql -h localhost -U aazriel -d postgres -c "CREATE USER scraper WITH PASSWORD 'scraper_password';"
psql -h localhost -U aazriel -d postgres -c "CREATE DATABASE scraper OWNER scraper;"

# Run migrations
DATABASE_URL="postgres://scraper:scraper_password@localhost:5432/scraper?sslmode=disable" make migrate-up
```

### 2. Deploy Services

```bash
# Start all services
docker-compose -f docker-compose.local.yml up -d

# Check service status
docker-compose -f docker-compose.local.yml ps
```

### 3. Verify Deployment

```bash
# Check API Gateway health
curl http://localhost:8082/health

# Check URL Manager logs
docker-compose -f docker-compose.local.yml logs url-manager
```

## 📖 Usage Guide

### API Endpoints

#### Health Check
```bash
curl http://localhost:8082/health
```

#### URL Management

**Create a URL for scraping:**
```bash
curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "frequency": "hourly",
    "description": "Example website for scraping"
  }'
```

**Get all URLs:**
```bash
curl http://localhost:8082/api/v1/urls
```

**Get URL by ID:**
```bash
curl http://localhost:8082/api/v1/urls/{url_id}
```

**Update URL:**
```bash
curl -X PUT http://localhost:8082/api/v1/urls/{url_id} \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://updated-example.com",
    "frequency": "daily",
    "description": "Updated description"
  }'
```

**Delete URL:**
```bash
curl -X DELETE http://localhost:8082/api/v1/urls/{url_id}
```

**Trigger scraping for a URL:**
```bash
curl -X POST http://localhost:8082/api/v1/urls/{url_id}/trigger
```

**Bulk trigger scraping:**
```bash
curl -X POST http://localhost:8082/api/v1/urls/bulk-trigger \
  -H "Content-Type: application/json" \
  -d '{
    "url_ids": ["url_id_1", "url_id_2", "url_id_3"]
  }'
```

### Data Endpoints

**Get scraping data:**
```bash
curl http://localhost:8082/api/v1/data
```

**Get data by URL ID:**
```bash
curl http://localhost:8082/api/v1/data/url/{url_id}
```

### Metrics Endpoints

**Get service metrics:**
```bash
curl http://localhost:8082/api/v1/metrics
```

**Get scraping statistics:**
```bash
curl http://localhost:8082/api/v1/metrics/scraping
```

## 🔍 Monitoring & Debugging

### View Service Logs

```bash
# All services
docker-compose -f docker-compose.local.yml logs -f

# Specific service
docker-compose -f docker-compose.local.yml logs -f url-manager
docker-compose -f docker-compose.local.yml logs -f api-gateway

# Using Docker directly
docker logs -f scraping_url_manager
docker logs -f scraping_api_gateway
```

### Kafka Monitoring

**Kafka UI (Web Interface):**
- Open [http://localhost:8080](http://localhost:8080)
- View topics, messages, and consumer groups
- Monitor message flow between services

**Kafka CLI:**
```bash
# List topics
docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list

# View messages in a topic
docker exec scraping_kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic scraping-tasks --from-beginning
```

### Database Queries

```bash
# Connect to database
psql -h localhost -U scraper -d scraper

# View URLs
SELECT * FROM urls;

# View scraping tasks
SELECT * FROM scraping_tasks;

# View recent scraping activity
SELECT url_id, status, created_at, updated_at 
FROM scraping_tasks 
ORDER BY created_at DESC 
LIMIT 10;
```

## 🛠️ Development

### Project Structure

```
go_scraping_project/
├── cmd/                    # Application entrypoints
│   ├── api-gateway/       # API Gateway service
│   └── url-manager/       # URL Manager service
├── internal/              # Private application code
│   ├── api-gateway/       # API Gateway logic
│   ├── url-manager/       # URL Manager logic
│   ├── database/          # Database queries (sqlc)
│   ├── domain/            # Domain models and interfaces
│   └── config/            # Configuration management
├── pkg/                   # Public packages
│   ├── database/          # Database connection
│   ├── kafka/             # Kafka producer/consumer
│   └── observability/     # Logging and metrics
├── sql/                   # Database migrations
│   └── schema/            # SQL schema files
├── configs/               # Configuration files
└── docker-compose.local.yml
```

### Common Commands

```bash
# Build services
make build-service SERVICE=url-manager
make build-service SERVICE=api-gateway

# Run services locally
make run-service SERVICE=url-manager
make run-service SERVICE=api-gateway

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Database migrations
make migrate-up
make migrate-down

# Generate sqlc code
make sqlc-generate
```

### Configuration

The application uses a hierarchical configuration system:

1. **Defaults** (hardcoded in `internal/config/config.go`)
2. **Config file** (`configs/config.yaml`)
3. **Environment variables** (with `SCRAPING_` prefix)

**Environment Variables:**
```bash
# Database
SCRAPING_DATABASE_HOST=host.docker.internal
SCRAPING_DATABASE_PORT=5432
SCRAPING_DATABASE_USER=scraper
SCRAPING_DATABASE_PASSWORD=scraper_password
SCRAPING_DATABASE_NAME=scraper

# Kafka
SCRAPING_KAFKA_BROKERS=kafka:29092
SCRAPING_KAFKA_GROUP_ID=url-manager-group

# Server
SCRAPING_SERVER_PORT=8080
SCRAPING_LOG_LEVEL=info
```

## 📊 Example Workflows

### 1. Basic URL Scraping Setup

```bash
# 1. Create a URL for scraping
curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://news.ycombinator.com",
    "frequency": "hourly",
    "description": "Hacker News homepage"
  }'

# 2. Check the URL was created
curl http://localhost:8082/api/v1/urls

# 3. Trigger immediate scraping
curl -X POST http://localhost:8082/api/v1/urls/{url_id}/trigger

# 4. Monitor Kafka for messages
docker exec scraping_kafka kafka-console-consumer \
  --bootstrap-server localhost:9092 \
  --topic scraping-tasks \
  --from-beginning
```

### 2. Bulk URL Management

```bash
# 1. Create multiple URLs
curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://github.com/trending",
    "frequency": "daily",
    "description": "GitHub trending repositories"
  }'

curl -X POST http://localhost:8082/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://stackoverflow.com/questions",
    "frequency": "daily",
    "description": "Stack Overflow questions"
  }'

# 2. Get all URLs
curl http://localhost:8082/api/v1/urls

# 3. Trigger scraping for all URLs
curl -X POST http://localhost:8082/api/v1/urls/bulk-trigger \
  -H "Content-Type: application/json" \
  -d '{
    "url_ids": ["url_id_1", "url_id_2"]
  }'
```

### 3. Monitoring and Debugging

```bash
# 1. Check service health
curl http://localhost:8082/health

# 2. View service logs
docker-compose -f docker-compose.local.yml logs -f url-manager

# 3. Check Kafka topics
docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list

# 4. Monitor Kafka UI
# Open http://localhost:8080 in your browser

# 5. Check database
psql -h localhost -U scraper -d scraper -c "SELECT * FROM urls;"
```

## 🔧 Troubleshooting

### Common Issues

**1. Database Connection Failed**
```bash
# Check if PostgreSQL is running
pg_isready -h localhost -p 5432

# Verify database and user exist
psql -h localhost -U scraper -d scraper -c "SELECT 1;"
```

**2. Kafka Connection Issues**
```bash
# Check Kafka health
docker-compose -f docker-compose.local.yml logs kafka

# Verify Kafka topics
docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list
```

**3. Service Won't Start**
```bash
# Check service logs
docker-compose -f docker-compose.local.yml logs -f url-manager

# Verify configuration
docker exec scraping_url_manager env | grep SCRAPING
```

**4. Build Failures**
```bash
# Update dependencies
go mod tidy

# Rebuild
docker-compose -f docker-compose.local.yml build url-manager
```

### Log Levels

Set log level via environment variable:
```bash
SCRAPING_LOG_LEVEL=debug  # More verbose logging
SCRAPING_LOG_LEVEL=info   # Default level
SCRAPING_LOG_LEVEL=warn   # Warnings and errors only
SCRAPING_LOG_LEVEL=error  # Errors only
```

## 🚀 Production Deployment

For production deployment, consider:

1. **Security**: Use proper secrets management
2. **Monitoring**: Add Prometheus/Grafana
3. **Scaling**: Use Kubernetes or Docker Swarm
4. **Backup**: Set up database backups
5. **SSL**: Use HTTPS for API endpoints

## 📝 API Reference

### URL Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/urls` | Create a new URL |
| GET | `/api/v1/urls` | List all URLs |
| GET | `/api/v1/urls/{id}` | Get URL by ID |
| PUT | `/api/v1/urls/{id}` | Update URL |
| DELETE | `/api/v1/urls/{id}` | Delete URL |
| POST | `/api/v1/urls/{id}/trigger` | Trigger scraping |
| POST | `/api/v1/urls/bulk-trigger` | Bulk trigger scraping |

### Data Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/data` | Get all scraping data |
| GET | `/api/v1/data/url/{id}` | Get data by URL ID |

### System

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/metrics` | Service metrics |
| GET | `/api/v1/metrics/scraping` | Scraping statistics |

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details. 