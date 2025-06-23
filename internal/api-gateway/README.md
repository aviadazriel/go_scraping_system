# API Gateway Service

This directory contains the API Gateway service implementation, which serves as the centralized entry point for all external client requests.

## Structure

```
internal/api-gateway/
├── handlers/           # HTTP request handlers
│   ├── url_handlers.go      # URL management endpoints
│   ├── data_handlers.go     # Data retrieval endpoints
│   ├── metrics_handlers.go  # Metrics and monitoring endpoints
│   ├── admin_handlers.go    # Admin and system management endpoints
│   ├── health_handlers.go   # Health check endpoints
│   ├── middleware.go        # HTTP middleware (logging, CORS, etc.)
│   └── router.go           # Route configuration and setup
└── README.md           # This file
```

## Architecture

The API Gateway follows a modular design pattern where:

- **Handlers**: Each handler file contains related HTTP endpoints grouped by functionality
- **Middleware**: Centralized middleware for cross-cutting concerns
- **Router**: Centralized route configuration that brings all handlers together

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication
Currently, the API Gateway does not require authentication. Future versions will support JWT-based authentication.

### Common Response Format
All endpoints return JSON responses with the following structure:
```json
{
  "data": {},           // Response data (varies by endpoint)
  "message": "string",  // Success/error message
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Error Responses
```json
{
  "error": "Error description",
  "code": "ERROR_CODE",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Handler Organization

### URL Handlers (`url_handlers.go`)

#### `POST /api/v1/urls` - Create a new URL for scraping
**Purpose**: Registers a new URL to be scraped with specified configuration.

**Request Body**:
```json
{
  "url": "https://example.com",
  "frequency": "1h",
  "parser_config": {
    "selectors": {
      "title": "h1",
      "content": ".content"
    }
  },
  "user_agent": "CustomBot/1.0",
  "timeout": 30,
  "rate_limit": 100,
  "max_retries": 3
}
```

**Response**:
```json
{
  "id": "url-123",
  "url": "https://example.com",
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Validation**:
- `url`: Required, must be a valid URL
- `frequency`: Required, format like "1h", "30m", "1d"
- `timeout`: Optional, seconds (default: 30)
- `rate_limit`: Optional, requests per minute (default: 60)
- `max_retries`: Optional, number of retries (default: 3)

#### `GET /api/v1/urls` - List all URLs
**Purpose**: Retrieves a paginated list of all registered URLs.

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page, max 100 (default: 20)

**Response**:
```json
{
  "urls": [
    {
      "id": "url-123",
      "url": "https://example.com",
      "frequency": "1h",
      "status": "active",
      "last_scraped_at": "2024-01-01T01:00:00Z",
      "next_scrape_at": "2024-01-01T02:00:00Z",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

#### `GET /api/v1/urls/{id}` - Get a specific URL
**Purpose**: Retrieves detailed information about a specific URL.

**Response**:
```json
{
  "id": "url-123",
  "url": "https://example.com",
  "frequency": "1h",
  "status": "active",
  "parser_config": {
    "selectors": {
      "title": "h1",
      "content": ".content"
    }
  },
  "user_agent": "CustomBot/1.0",
  "timeout": 30,
  "rate_limit": 100,
  "max_retries": 3,
  "last_scraped_at": "2024-01-01T01:00:00Z",
  "next_scrape_at": "2024-01-01T02:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:30:00Z"
}
```

#### `PUT /api/v1/urls/{id}` - Update a URL
**Purpose**: Updates configuration for an existing URL.

**Request Body** (all fields optional):
```json
{
  "frequency": "2h",
  "parser_config": {
    "selectors": {
      "title": "h1.title",
      "content": ".main-content"
    }
  },
  "user_agent": "UpdatedBot/2.0",
  "timeout": 45,
  "rate_limit": 150,
  "max_retries": 5
}
```

**Response**:
```json
{
  "message": "URL updated successfully"
}
```

#### `DELETE /api/v1/urls/{id}` - Delete a URL
**Purpose**: Removes a URL from the scraping schedule.

**Response**:
```json
{
  "message": "URL deleted successfully"
}
```

#### `POST /api/v1/urls/{id}/scrape` - Trigger immediate scraping
**Purpose**: Manually triggers scraping for a specific URL, bypassing the schedule.

**Response**:
```json
{
  "message": "Scrape triggered successfully"
}
```

#### `GET /api/v1/urls/{id}/status` - Get URL status
**Purpose**: Retrieves current status and scheduling information for a URL.

**Response**:
```json
{
  "id": "url-123",
  "status": "active",
  "last_scraped_at": "2024-01-01T01:00:00Z",
  "next_scrape_at": "2024-01-01T02:00:00Z",
  "retry_count": 0,
  "max_retries": 3,
  "last_error": null
}
```

### Data Handlers (`data_handlers.go`)

#### `GET /api/v1/data` - List scraped data
**Purpose**: Retrieves a paginated list of scraped and parsed data.

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page, max 100 (default: 20)
- `schema`: Filter by data schema
- `url_id`: Filter by URL ID

**Response**:
```json
{
  "data": [
    {
      "id": "data-123",
      "url_id": "url-123",
      "schema": "article",
      "data": {
        "title": "Example Article",
        "content": "Article content...",
        "author": "John Doe",
        "published_at": "2024-01-01T00:00:00Z"
      },
      "scraped_at": "2024-01-01T01:00:00Z",
      "created_at": "2024-01-01T01:00:00Z"
    }
  ],
  "total": 1000,
  "page": 1,
  "limit": 20
}
```

#### `GET /api/v1/data/{url_id}` - Get data for a specific URL
**Purpose**: Retrieves all scraped data for a specific URL.

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page, max 100 (default: 20)

**Response**: Same format as list data, but filtered by URL ID.

#### `GET /api/v1/data/export` - Export data in various formats
**Purpose**: Exports scraped data in JSON, CSV, or XML format.

**Query Parameters**:
- `format`: Export format (`json`, `csv`, `xml`) - default: `json`
- `url_id`: Filter by URL ID (can be multiple)
- `schema`: Filter by data schema
- `from`: Start date (ISO 8601)
- `to`: End date (ISO 8601)

**Response**:
```json
{
  "format": "json",
  "count": 1000,
  "data": [...],
  "exported_at": "2024-01-01T01:00:00Z"
}
```

### Metrics Handlers (`metrics_handlers.go`)

#### `GET /api/v1/metrics/urls/{id}` - Get metrics for a specific URL
**Purpose**: Retrieves performance and success metrics for a specific URL.

**Query Parameters**:
- `period`: Time period (`1h`, `24h`, `7d`, `30d`) - default: `24h`
- `include_time_series`: Include time series data (`true`/`false`) - default: `false`

**Response**:
```json
{
  "url_id": "url-123",
  "total_requests": 100,
  "success_rate": 95.5,
  "avg_response_time": 250,
  "last_scraped_at": "2024-01-01T01:00:00Z",
  "status_counts": {
    "200": 95,
    "404": 3,
    "500": 2
  },
  "error_counts": {
    "timeout": 2,
    "connection": 1,
    "parse_error": 0
  },
  "time_series": [
    {
      "timestamp": "2024-01-01T00:00:00Z",
      "response_time": 200,
      "status_code": 200,
      "success": true
    }
  ]
}
```

#### `GET /api/v1/metrics/system` - Get system-wide metrics
**Purpose**: Retrieves overall system performance and health metrics.

**Query Parameters**:
- `period`: Time period (`1h`, `24h`, `7d`, `30d`) - default: `24h`

**Response**:
```json
{
  "total_urls": 50,
  "active_urls": 45,
  "total_scraped_data": 1000,
  "total_parsed_data": 950,
  "success_rate": 95.0,
  "avg_response_time": 275,
  "dead_letter_count": 5,
  "retry_count": 12,
  "last_updated": "2024-01-01T01:00:00Z"
}
```

### Admin Handlers (`admin_handlers.go`)

#### `GET /api/v1/admin/dead-letter` - List dead letter messages
**Purpose**: Retrieves messages that failed processing and are in the dead letter queue.

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page, max 100 (default: 20)
- `topic`: Filter by Kafka topic
- `status`: Filter by status (`pending`, `retrying`, `failed`)

**Response**:
```json
{
  "messages": [
    {
      "id": "msg-123",
      "topic": "scraping-requests",
      "partition": 0,
      "offset": 12345,
      "error": "Connection timeout",
      "retry_count": 3,
      "max_retries": 3,
      "next_retry_at": null,
      "failed_at": "2024-01-01T01:00:00Z",
      "original_message": {
        "url_id": "url-123",
        "url": "https://example.com"
      }
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 20
}
```

#### `POST /api/v1/admin/dead-letter/bulk-retry` - Bulk retry failed messages
**Purpose**: Retries multiple failed messages at once.

**Request Body**:
```json
{
  "topic": "scraping-requests",
  "message_ids": ["msg-123", "msg-124"],
  "force_retry": false
}
```

**Response**:
```json
{
  "message": "Bulk retry initiated",
  "retried": 2,
  "failed": 0
}
```

#### `POST /api/v1/admin/dead-letter/{id}/retry` - Retry a specific message
**Purpose**: Retries a specific failed message.

**Request Body** (optional):
```json
{
  "force_retry": true
}
```

**Response**:
```json
{
  "message": "Message retry initiated successfully"
}
```

#### `DELETE /api/v1/admin/dead-letter/{id}` - Delete a dead letter message
**Purpose**: Permanently removes a message from the dead letter queue.

**Response**:
```json
{
  "message": "Message deleted successfully"
}
```

#### `GET /api/v1/admin/health` - Get system health status
**Purpose**: Retrieves comprehensive health status of all system components.

**Response**:
```json
{
  "status": "healthy",
  "services": {
    "database": {
      "status": "healthy",
      "latency": 5
    },
    "kafka": {
      "status": "healthy",
      "brokers": 1
    },
    "scraper": {
      "status": "healthy",
      "active_workers": 3
    }
  },
  "timestamp": "2024-01-01T01:00:00Z"
}
```

### Health Handlers (`health_handlers.go`)

#### `GET /health` - Basic health check
**Purpose**: Simple health check endpoint for load balancers and monitoring.

**Response**:
```json
{
  "status": "healthy",
  "service": "api-gateway",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "checks": {
    "database": "healthy",
    "kafka": "healthy"
  }
}
```

#### `GET /ready` - Readiness probe
**Purpose**: Kubernetes readiness probe to check if the service is ready to receive traffic.

**Response**:
```json
{
  "status": "ready",
  "service": "api-gateway",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

#### `GET /live` - Liveness probe
**Purpose**: Kubernetes liveness probe to check if the service is alive and responsive.

**Response**:
```json
{
  "status": "alive",
  "service": "api-gateway",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Middleware

The API Gateway includes several middleware components:

- **Logging Middleware**: Structured logging for all HTTP requests
- **CORS Middleware**: Cross-Origin Resource Sharing support
- **Recovery Middleware**: Panic recovery and error handling
- **Auth Middleware**: Authentication (placeholder for future implementation)
- **Rate Limit Middleware**: Rate limiting (placeholder for future implementation)
- **Request ID Middleware**: Request tracing (placeholder for future implementation)

## Benefits of This Structure

1. **Service-Specific Organization**: All API Gateway code is contained within its own directory
2. **Modular Handlers**: Each handler file focuses on a specific domain area
3. **Easy Maintenance**: Clear separation of concerns makes the code easier to maintain
4. **Scalability**: New handlers can be added without affecting existing ones
5. **Testability**: Each handler can be tested independently
6. **Future-Proof**: Structure supports adding more services without confusion

## Usage

The API Gateway is started from `cmd/api-gateway/main.go` and uses the handlers from this package to handle all incoming HTTP requests.

## Future Enhancements

- Add authentication and authorization middleware
- Implement rate limiting
- Add request ID generation and tracing
- Add metrics collection middleware
- Implement circuit breaker patterns
- Add API versioning support 