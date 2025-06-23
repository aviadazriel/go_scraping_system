# API Gateway Service

The API Gateway serves as the single entry point for external clients to interact with the scraping microservices. It provides a RESTful API for managing URLs, retrieving scraped data, and monitoring system health.

## Features

- **URL Management**: Create, read, update, and delete URLs for scraping
- **Data Retrieval**: Access parsed data with filtering and pagination
- **Real-time Status**: Check scraping status and trigger immediate scraping
- **Metrics**: Monitor system performance and URL-specific metrics
- **Admin Operations**: Manage dead letter queue and retry failed operations
- **Health Monitoring**: Built-in health checks and observability

## API Endpoints

### Health Check
- `GET /health` - Service health status

### URL Management
- `POST /api/v1/urls` - Create a new URL for scraping
- `GET /api/v1/urls` - List all URLs with pagination
- `GET /api/v1/urls/{id}` - Get specific URL details
- `PUT /api/v1/urls/{id}` - Update URL configuration
- `DELETE /api/v1/urls/{id}` - Delete a URL
- `POST /api/v1/urls/{id}/scrape` - Trigger immediate scraping
- `GET /api/v1/urls/{id}/status` - Get URL scraping status

### Data Retrieval
- `GET /api/v1/data` - List parsed data with filtering
- `GET /api/v1/data/{url_id}` - Get data for specific URL
- `GET /api/v1/data/export` - Export data in various formats

### Metrics
- `GET /api/v1/metrics/urls/{id}` - Get URL-specific metrics
- `GET /api/v1/metrics/system` - Get system-wide metrics

### Admin Operations
- `GET /api/v1/admin/dead-letter` - List dead letter messages
- `POST /api/v1/admin/dead-letter/{id}/retry` - Retry failed message
- `DELETE /api/v1/admin/dead-letter/{id}` - Delete dead letter message

## Request/Response Examples

### Create URL
```bash
curl -X POST http://localhost:8080/api/v1/urls \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "frequency": "1h",
    "parser_config": {
      "title_selector": "h1",
      "content_selector": ".content"
    },
    "user_agent": "Mozilla/5.0 (compatible; Scraper/1.0)",
    "timeout": 30,
    "rate_limit": 1,
    "max_retries": 3
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "url": "https://example.com",
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z"
}
```

### List URLs
```bash
curl "http://localhost:8080/api/v1/urls?page=1&limit=20"
```

Response:
```json
{
  "urls": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "url": "https://example.com",
      "frequency": "1h",
      "status": "completed",
      "last_scraped_at": "2024-01-01T01:00:00Z",
      "next_scrape_at": "2024-01-01T02:00:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1,
  "page": 1,
  "limit": 20
}
```

### Get URL Status
```bash
curl http://localhost:8080/api/v1/urls/550e8400-e29b-41d4-a716-446655440000/status
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "last_scraped_at": "2024-01-01T01:00:00Z",
  "next_scrape_at": "2024-01-01T02:00:00Z",
  "retry_count": 0,
  "max_retries": 3
}
```

## Configuration

The service can be configured using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `scraping_db` | Database name |
| `DB_USER` | `scraping_user` | Database user |
| `DB_PASSWORD` | `scraping_password` | Database password |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses |
| `LOG_LEVEL` | `info` | Logging level |

## Running the Service

### Development
```bash
# Build the service
go build -o bin/api-gateway ./cmd/api-gateway

# Run the service
./bin/api-gateway
```

### Docker
```bash
# Build the Docker image
docker build -f cmd/api-gateway/Dockerfile -t api-gateway .

# Run the container
docker run -p 8080:8080 api-gateway
```

### Docker Compose
```bash
# Start all services including API Gateway
docker-compose up api-gateway
```

## Health Check

The service provides a health check endpoint at `/health` that returns:

```json
{
  "status": "healthy",
  "service": "api-gateway",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Error Handling

The API Gateway returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses include a JSON object with an `error` field:

```json
{
  "error": "URL not found"
}
```

## Middleware

The service includes several middleware components:

- **Logging**: Structured logging for all HTTP requests
- **CORS**: Cross-origin resource sharing support
- **Recovery**: Panic recovery with proper error responses
- **Authentication**: (Future) JWT-based authentication
- **Rate Limiting**: (Future) Request rate limiting

## Monitoring

The service integrates with:

- **Prometheus**: Metrics collection
- **Jaeger**: Distributed tracing
- **Structured Logging**: JSON-formatted logs

## Development

### Project Structure
```
cmd/api-gateway/
├── main.go              # Service entry point
├── Dockerfile           # Docker configuration
└── README.md           # This file

internal/handlers/
├── url_handlers.go     # URL management handlers
└── data_handlers.go    # Data retrieval handlers
```

### Adding New Endpoints

1. Create a new handler function in the appropriate handler file
2. Add the route in `main.go` in the `createRouter` function
3. Add request/response types if needed
4. Update this README with the new endpoint documentation

### Testing

```bash
# Run unit tests
go test ./internal/handlers/...

# Run integration tests
go test ./cmd/api-gateway/...
```

## Dependencies

- **Gorilla Mux**: HTTP router and URL matching
- **Logrus**: Structured logging
- **Viper**: Configuration management
- **OpenTelemetry**: Observability and tracing 