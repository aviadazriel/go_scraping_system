# Go Scraping Microservice System

A distributed web scraping system built with Go microservices, Kafka, and PostgreSQL.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Gateway   │    │  URL Manager    │    │   Scraper       │
│   (Web Service) │    │  (Background    │    │   Service       │
│                 │    │   Service)      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Kafka Topics  │    │   PostgreSQL    │    │   HTML Storage  │
│                 │    │   (URLs/Data)   │    │   (File System) │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Parser        │    │   Data Storage  │    │   Dead Letter   │
│   Service       │    │   Service       │    │   Handler       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Services

### 1. URL Manager Service (Background)
- **Purpose**: **Background service** that schedules and distributes web scraping tasks
- **Responsibilities**:
  - **Periodically scans database** for URLs due for scraping (every 30 seconds)
  - **Creates scraping tasks** and sends them to Kafka topics
  - **Updates URL scheduling** (next scrape time, status, retry counts)
  - **Task distribution** to scraper services via Kafka
  - **No HTTP endpoints** - pure background processing

### 2. API Gateway Service (Web)
- **Purpose**: Single entry point for external clients
- **Responsibilities**:
  - **CRUD operations** for URLs (create, read, update, delete)
  - **Manual trigger endpoints** for immediate scraping
  - **Data retrieval** and export functionality
  - **Authentication and authorization**
  - **Rate limiting and request transformation**

### 3. Scraper Service
- **Purpose**: Performs HTTP requests and saves HTML content
- **Responsibilities**:
  - Consume scraping tasks from Kafka
  - HTTP requests with retry logic and rate limiting
  - HTML file storage
  - Status updates to database
  - Error handling and dead letter queue

### 4. Parser Service
- **Purpose**: Extracts structured data from HTML files
- **Responsibilities**:
  - HTML parsing and data extraction
  - Data validation and transformation
  - Structured data storage
  - Error handling for malformed HTML

### 5. Data Storage Service
- **Purpose**: Manages parsed data storage and retrieval
- **Responsibilities**:
  - Store parsed data in database
  - Data querying and filtering
  - Data export functionality
  - Data cleanup and archiving

## Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL 15+
- **Message Queue**: Apache Kafka
- **API**: gRPC + REST
- **Observability**: OpenTelemetry, Prometheus, Grafana
- **Containerization**: Docker & Docker Compose
- **Orchestration**: Kubernetes (optional)

## Project Structure

```
go_scraping_project/
├── cmd/
│   ├── api-gateway/
│   ├── url-manager/
│   ├── scraper/
│   ├── parser/
│   └── data-storage/
├── internal/
│   ├── domain/
│   ├── services/
│   ├── repositories/
│   ├── handlers/
│   └── config/
├── pkg/
│   ├── kafka/
│   ├── database/
│   ├── observability/
│   └── utils/
├── api/
│   ├── proto/
│   └── openapi/
├── configs/
├── deployments/
├── scripts/
└── docs/
```

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+
- Apache Kafka

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go_scraping_project
   ```

2. **Start infrastructure services**
   ```bash
   docker-compose up -d postgres kafka zookeeper
   ```

3. **Run database migrations**
   ```bash
   make migrate
   ```

4. **Start all services**
   ```bash
   make run-all
   ```

5. **Add a URL to scrape**
   ```bash
   curl -X POST http://localhost:8080/api/v1/urls \
     -H "Content-Type: application/json" \
     -d '{
       "url": "https://example.com",
       "frequency": "1h",
       "parser_config": {
         "title_selector": "h1",
         "content_selector": ".content"
       }
     }'
   ```

## Configuration

Each service can be configured via environment variables or configuration files:

```bash
# Database
DATABASE_URL=postgres://user:password@localhost:5432/scraping_db

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=scraper-group

# Service specific
SERVICE_PORT=8080
LOG_LEVEL=info
```

## API Documentation

### URL Management
- `POST /api/v1/urls` - Add new URL to scrape
- `GET /api/v1/urls` - List all URLs
- `GET /api/v1/urls/{id}` - Get URL details
- `PUT /api/v1/urls/{id}` - Update URL
- `DELETE /api/v1/urls/{id}` - Delete URL

### Scraping Status
- `GET /api/v1/urls/{id}/status` - Get scraping status
- `POST /api/v1/urls/{id}/scrape` - Trigger immediate scrape

### Data Retrieval
- `GET /api/v1/data` - Get parsed data
- `GET /api/v1/data/{url_id}` - Get data for specific URL
- `GET /api/v1/data/export` - Export data

## Monitoring & Observability

### Metrics
- Scraping success/failure rates
- Response times
- Queue depths
- Database connection health

### Logging
- Structured JSON logging
- Request correlation IDs
- Error tracking

### Tracing
- Distributed tracing across services
- Performance monitoring
- Error analysis

## Development

### Running Tests
```bash
make test
make test-coverage
```

### Code Quality
```bash
make lint
make fmt
make vet
```

### Building
```bash
make build
make docker-build
```

## Deployment

### Docker Compose
```bash
docker-compose up -d
```

### Kubernetes
```bash
kubectl apply -f deployments/k8s/
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details 