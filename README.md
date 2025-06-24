# Go Scraping Project

A microservice-based web scraping system built with Go, featuring Kafka for event-driven architecture, PostgreSQL for data persistence, and Docker for containerization.

## 🏗️ Architecture

The system consists of multiple microservices:

- **API Gateway**: HTTP API for managing URLs and scraping tasks
- **URL Manager**: Background service for scheduling and triggering scraping tasks
- **Scraper Service**: (Coming soon) Service for actual web scraping
- **Parser Service**: (Coming soon) Service for parsing scraped HTML
- **Database Service**: PostgreSQL with sqlc for type-safe database operations
- **Kafka**: Event streaming platform for inter-service communication

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- PostgreSQL (running locally or in Docker)
- Make (optional, for using Makefile commands)

### Option 1: Automated Deployment (Recommended)

Use the automated deployment script that will set up everything for you:

```bash
./deploy.sh
```

This script will:
- Check your PostgreSQL connection
- Create the database and user if needed
- Run database migrations
- Build and start all services
- Verify all services are healthy

### Option 2: Manual Deployment

#### 1. Set up PostgreSQL

Make sure PostgreSQL is running and create the database:

```bash
# Create database and user
psql -U postgres -c "CREATE USER scraper WITH PASSWORD 'scraper_password';"
psql -U postgres -c "CREATE DATABASE scraper OWNER scraper;"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE scraper TO scraper;"
```

#### 2. Run Database Migrations

```bash
# Install goose if you haven't already
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/migrations postgres "host=localhost port=5432 user=scraper password=scraper_password dbname=scraper sslmode=disable" up
```

#### 3. Start Services

```bash
# Start all services (uses your local PostgreSQL)
docker-compose -f docker-compose.local.yml up -d --build

# Or start with PostgreSQL in Docker
docker-compose -f docker-compose.production.yml up -d --build
```

## 📋 Service Endpoints

### API Gateway (Port 8082)
- **Health Check**: `GET http://localhost:8082/health`
- **Create URL**: `POST http://localhost:8082/api/v1/urls`
- **List URLs**: `GET http://localhost:8082/api/v1/urls`
- **Get URL**: `GET http://localhost:8082/api/v1/urls/{id}`

### URL Manager (Port 8081)
- **Health Check**: `GET http://localhost:8081/health`
- **Trigger All URLs**: `POST http://localhost:8081/trigger/all`
- **Trigger Specific URL**: `POST http://localhost:8081/trigger/{id}`

### Kafka UI (Port 8080)
- **Web Interface**: `http://localhost:8080`

## 🧪 Testing the API

### Create a URL for scraping

```bash
curl -X POST http://localhost:8082/api/v1/urls \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://example.com",
    "frequency": "1h",
    "user_agent": "GoScrapingBot/1.0",
    "timeout": 30,
    "rate_limit": 1,
    "max_retries": 3
  }'
```

### List all URLs

```bash
curl http://localhost:8082/api/v1/urls
```

### Check service health

```bash
curl http://localhost:8082/health
curl http://localhost:8081/health
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
│   ├── domain/            # Domain models and interfaces
│   ├── database/          # Database queries (sqlc-generated)
│   └── config/            # Configuration management
├── pkg/                   # Public packages
│   ├── kafka/             # Kafka producer/consumer
│   ├── database/          # Database connection
│   └── observability/     # Logging and monitoring
├── sql/                   # Database migrations and queries
│   ├── migrations/        # Goose migration files
│   └── queries/           # sqlc query files
├── configs/               # Configuration files
├── docker-compose.local.yml    # Local development setup
├── docker-compose.production.yml # Production setup
└── deploy.sh              # Automated deployment script
```

### Building Services

```bash
# Build all services
make build

# Build specific service
make build-api-gateway
make build-url-manager
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Database Operations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Generate sqlc code
make sqlc-generate
```

## 🔧 Configuration

The services use environment variables for configuration. Key settings:

### Database
- `DB_HOST`: PostgreSQL host (default: localhost)
- `DB_PORT`: PostgreSQL port (default: 5432)
- `DB_USER`: Database user (default: scraper)
- `DB_PASSWORD`: Database password (default: scraper_password)
- `DB_NAME`: Database name (default: scraper)

### Kafka
- `KAFKA_BROKERS`: Kafka broker addresses (default: localhost:9092)
- `KAFKA_GROUP_ID`: Consumer group ID
- `KAFKA_RETRY_MAX_ATTEMPTS`: Maximum retry attempts
- `KAFKA_RETRY_BACKOFF`: Retry backoff duration

### Server
- `SERVER_PORT`: HTTP server port (default: 8080)
- `SERVER_READ_TIMEOUT`: Read timeout
- `SERVER_WRITE_TIMEOUT`: Write timeout
- `SERVER_IDLE_TIMEOUT`: Idle timeout

## 📊 Monitoring

### Health Checks
All services expose health check endpoints:
- API Gateway: `http://localhost:8082/health`
- URL Manager: `http://localhost:8081/health`

### Logs
View service logs:
```bash
# All services
docker-compose -f docker-compose.local.yml logs -f

# Specific service
docker-compose -f docker-compose.local.yml logs -f url-manager
```

### Kafka UI
Monitor Kafka topics and messages at `http://localhost:8080`

## 🚀 Production Deployment

For production deployment, use the production Docker Compose file:

```bash
docker-compose -f docker-compose.production.yml up -d --build
```

This includes:
- PostgreSQL in Docker
- All microservices
- Proper networking and volumes
- Health checks and restart policies

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details. 