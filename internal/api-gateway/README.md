# API Gateway - Internal Structure

This directory contains the internal implementation of the API Gateway service, organized in a modular and maintainable structure.

## Directory Structure

```
internal/api-gateway/
├── handlers/           # HTTP request handlers (minimal, just route setup)
│   ├── router.go       # Route configuration and setup
│   ├── middleware.go   # HTTP middleware (logging, CORS, recovery)
│   ├── health_handlers.go # Health check endpoints (simple, no models needed)
│   ├── url_handlers.go     # Placeholder (functionality moved to types)
│   ├── data_handlers.go    # Placeholder (functionality moved to types)
│   ├── metrics_handlers.go # Placeholder (functionality moved to types)
│   └── admin_handlers.go   # Placeholder (functionality moved to types)
├── models/             # Data models and structs
│   ├── requests.go     # All request structs
│   ├── responses.go    # All response structs
│   ├── common.go       # Shared types and errors
│   └── config.go       # Configuration structs
└── types/              # Handler type definitions and implementations
    ├── handlers.go     # Router struct and dependencies
    ├── url_handler.go  # URLHandler struct and implementation
    ├── data_handler.go # DataHandler struct and implementation
    ├── metrics_handler.go # MetricsHandler struct and implementation
    └── admin_handler.go # AdminHandler struct and implementation
```

## Architecture Overview

### Separation of Concerns

The API Gateway follows a clean separation of concerns:

1. **Handlers** (`handlers/`): Contain route setup, middleware, and simple health checks
2. **Models** (`models/`): Define data structures for requests, responses, and shared types
3. **Types** (`types/`): Define handler structs and their complete implementations

### Benefits of This Structure

#### 1. **Improved Maintainability**
- Each file has a single responsibility
- Easy to locate and modify specific functionality
- Clear boundaries between different concerns

#### 2. **Better Reusability**
- Models can be imported and used across different handlers
- Common types are centralized and consistent
- Reduced code duplication

#### 3. **Enhanced Testability**
- Models can be tested independently
- Handlers can be mocked more easily
- Clear dependencies make unit testing straightforward

#### 4. **Scalability**
- Easy to add new handlers without cluttering existing files
- Models can be extended without affecting handler logic
- Clear structure for new team members

## File Descriptions

### Handlers (`handlers/`)

- **`router.go`**: Route configuration, handler initialization, and route setup functions
- **`middleware.go`**: HTTP middleware for logging, CORS, and error handling
- **`health_handlers.go`**: Health check endpoints (health, ready, live)
- **`*_handlers.go`**: Placeholder files with comments explaining the refactoring

### Models (`models/`)

- **`requests.go`**: All request structs used by handlers
  - `CreateURLRequest`
  - `UpdateURLRequest`
  - `ExportDataRequest`
  - `BulkRetryRequest`

- **`responses.go`**: All response structs returned by handlers
  - `CreateURLResponse`
  - `ListURLsResponse`
  - `URLMetricsResponse`
  - `SystemMetricsResponse`
  - `HealthResponse`
  - `DeadLetterMessageResponse`

- **`common.go`**: Shared types and utilities
  - `ValidationError`
  - `responseWriter` (for middleware)

- **`config.go`**: Configuration-related structs
  - `ParserConfig`

### Types (`types/`)

- **`handlers.go`**: Router struct definition with handler dependencies
  - `Router` (contains references to all handler instances)

- **`url_handler.go`**: URLHandler struct definition and complete implementation
  - `URLHandler` struct
  - `NewURLHandler` constructor
  - `CreateURL`
  - `ListURLs`
  - `GetURL`
  - `UpdateURL`
  - `DeleteURL`
  - `TriggerScrape`
  - `GetURLStatus`

- **`data_handler.go`**: DataHandler struct definition and complete implementation
  - `DataHandler` struct
  - `NewDataHandler` constructor
  - `ListData`
  - `GetDataByURL`
  - `ExportData`

- **`metrics_handler.go`**: MetricsHandler struct definition and complete implementation
  - `MetricsHandler` struct
  - `NewMetricsHandler` constructor
  - `GetURLMetrics`
  - `GetSystemMetrics`

- **`admin_handler.go`**: AdminHandler struct definition and complete implementation
  - `AdminHandler` struct
  - `NewAdminHandler` constructor
  - `ListDeadLetterMessages`
  - `RetryDeadLetterMessage`
  - `DeleteDeadLetterMessage`
  - `BulkRetryDeadLetterMessages`
  - `GetSystemHealth`

## Usage Examples

### Creating a New Handler

1. **Define the handler struct** in `types/new_handler.go`:
```go
type NewHandler struct {
    Logger *logrus.Logger
    DB     *database.Queries
}
```

2. **Define request/response models** in `models/requests.go` and `models/responses.go`:
```go
type NewRequest struct {
    Field string `json:"field"`
}

type NewResponse struct {
    Result string `json:"result"`
}
```

3. **Implement the handler** in `types/new_handler.go`:
```go
func NewNewHandler(logger *logrus.Logger, db *database.Queries) *NewHandler {
    return &NewHandler{
        Logger: logger,
        DB:     db,
    }
}

func (h *NewHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    var req models.NewRequest
    // ... implementation
}
```

4. **Register routes** in `handlers/router.go`:
```go
func setupNewRoutes(apiV1 *mux.Router, newHandler *types.NewHandler) {
    newRoutes := apiV1.PathPrefix("/new").Subrouter()
    newRoutes.HandleFunc("", newHandler.HandleRequest).Methods("POST")
}
```

### Importing Models

```go
import (
    "go_scraping_project/internal/api-gateway/models"
    "go_scraping_project/internal/api-gateway/types"
)

func someFunction() {
    req := models.CreateURLRequest{
        URL:       "https://example.com",
        Frequency: "1h",
    }
    
    handler := types.NewURLHandler(logger, db)
}
```

## Best Practices

1. **Keep handlers focused**: Each handler should handle one specific endpoint
2. **Use models consistently**: Always use the defined models for requests/responses
3. **Document structs**: Add comprehensive comments for all exported types
4. **Validate inputs**: Use the `ValidationError` type for consistent error handling
5. **Follow naming conventions**: Use descriptive names that indicate purpose

## Migration from Previous Structure

The previous structure had all structs and implementations defined within handler files. This new structure:

- ✅ Separates concerns clearly
- ✅ Improves code organization
- ✅ Makes the codebase more maintainable
- ✅ Reduces duplication
- ✅ Improves testability
- ✅ Provides clear separation between models, types, and handlers

All existing functionality remains the same, but the code is now better organized and more maintainable.

## API Endpoints

### URL Management
- `POST /api/v1/urls` - Create a new URL
- `GET /api/v1/urls` - List all URLs (with pagination)
- `GET /api/v1/urls/{id}` - Get specific URL details
- `PUT /api/v1/urls/{id}` - Update URL configuration
- `DELETE /api/v1/urls/{id}` - Delete a URL
- `POST /api/v1/urls/{id}/scrape` - Trigger manual scraping
- `GET /api/v1/urls/{id}/status` - Get URL status information

### Data Management
- `GET /api/v1/data` - List scraped data (with filtering and pagination)
- `GET /api/v1/data/{url_id}` - Get data for specific URL
- `GET /api/v1/data/export` - Export data in various formats

### Metrics
- `GET /api/v1/metrics/urls/{id}` - Get metrics for specific URL
- `GET /api/v1/metrics/system` - Get system-wide metrics

### Admin
- `GET /api/v1/admin/dead-letter` - List dead letter messages
- `POST /api/v1/admin/dead-letter/bulk-retry` - Bulk retry failed messages
- `POST /api/v1/admin/dead-letter/{id}/retry` - Retry specific message
- `DELETE /api/v1/admin/dead-letter/{id}` - Delete dead letter message
- `GET /api/v1/admin/health` - Get comprehensive system health

### Health Checks
- `GET /health` - Basic health check
- `GET /ready` - Readiness probe
- `GET /live` - Liveness probe 