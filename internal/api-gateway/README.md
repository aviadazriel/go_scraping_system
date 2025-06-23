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

## Handler Organization

### URL Handlers (`url_handlers.go`)
- `POST /api/v1/urls` - Create a new URL for scraping
- `GET /api/v1/urls` - List all URLs
- `GET /api/v1/urls/{id}` - Get a specific URL
- `PUT /api/v1/urls/{id}` - Update a URL
- `DELETE /api/v1/urls/{id}` - Delete a URL
- `POST /api/v1/urls/{id}/scrape` - Trigger immediate scraping
- `GET /api/v1/urls/{id}/status` - Get URL status

### Data Handlers (`data_handlers.go`)
- `GET /api/v1/data` - List scraped data
- `GET /api/v1/data/{url_id}` - Get data for a specific URL
- `GET /api/v1/data/export` - Export data in various formats

### Metrics Handlers (`metrics_handlers.go`)
- `GET /api/v1/metrics/urls/{id}` - Get metrics for a specific URL
- `GET /api/v1/metrics/system` - Get system-wide metrics

### Admin Handlers (`admin_handlers.go`)
- `GET /api/v1/admin/dead-letter` - List dead letter messages
- `POST /api/v1/admin/dead-letter/bulk-retry` - Bulk retry failed messages
- `POST /api/v1/admin/dead-letter/{id}/retry` - Retry a specific message
- `DELETE /api/v1/admin/dead-letter/{id}` - Delete a dead letter message
- `GET /api/v1/admin/health` - Get system health status

### Health Handlers (`health_handlers.go`)
- `GET /health` - Basic health check
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe

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