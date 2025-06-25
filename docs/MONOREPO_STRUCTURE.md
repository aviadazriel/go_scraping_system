# Multi-Project Monorepo Structure

This document explains the new multi-project monorepo structure for the Go Scraping Project.

## Overview

The project has been refactored from a single-project structure to a multi-project monorepo where each service has its own `go.mod`, Dockerfile, and Makefile while sharing common code through a shared package.

## Structure Comparison

### Before (Single-Project)
```
go_scraping_project/
├── cmd/
│   ├── api-gateway/
│   └── url-manager/
├── internal/
│   ├── api-gateway/
│   ├── url-manager/
│   └── database/
├── pkg/
│   └── utils/
├── go.mod
└── docker-compose.yml
```

### After (Multi-Project Monorepo)
```
go_scraping_project/
├── shared/                    # Shared packages
│   ├── utils/
│   ├── models/
│   ├── config/
│   └── go.mod
├── services/                  # Individual services
│   ├── api-gateway/
│   │   ├── handlers/
│   │   ├── database/
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── Dockerfile
│   │   └── Makefile
│   ├── url-manager/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── repositories/
│   │   ├── database/
│   │   ├── main.go
│   │   ├── go.mod
│   │   ├── Dockerfile
│   │   └── Makefile
│   └── [future services]
├── docker-compose.yml
└── Makefile
```

## Benefits

### 1. Service Independence
- Each service has its own dependencies and can evolve independently
- Services can be developed, tested, and deployed separately
- Different teams can work on different services without conflicts

### 2. Clear Boundaries
- Explicit separation between services
- Shared code is clearly identified in the `shared/` directory
- Each service has its own entry point and configuration

### 3. Scalability
- Easy to add new services without affecting existing ones
- Services can be split into separate repositories later if needed
- Independent versioning and release cycles

### 4. Development Experience
- Faster builds (only rebuild changed services)
- Independent testing and linting
- Service-specific tooling and configuration

## Shared Code

The `shared/` directory contains code that is used across multiple services:

### `shared/utils/`
- Time utilities
- Validation functions
- Common helper functions

### `shared/models/`
- Domain models (URL, ScrapingTask, etc.)
- Data structures used across services

### `shared/config/`
- Configuration structures
- Default configuration values

## Service Structure

Each service follows a consistent structure:

```
services/[service-name]/
├── handlers/          # HTTP handlers (if applicable)
├── services/          # Business logic
├── repositories/      # Data access layer
├── database/          # Database-specific code
├── main.go           # Service entry point
├── go.mod            # Service dependencies
├── Dockerfile        # Service container
└── Makefile          # Build commands
```

## Module Dependencies

### Shared Module
```go
module go_scraping_project/shared

go 1.21

require (
    github.com/google/uuid v1.4.0
)
```

### Service Modules
```go
module go_scraping_project/services/[service-name]

go 1.21

require (
    // Service-specific dependencies
    go_scraping_project/shared v0.0.0
)

replace go_scraping_project/shared => ../../shared
```

## Build System

### Top-Level Makefile
The root `Makefile` provides commands to work with all services:

```bash
# Build all services
make build-all

# Test all services
make test-all

# Run specific service
make api-gateway
make url-manager
```

### Service-Specific Makefiles
Each service has its own `Makefile` for service-specific operations:

```bash
# Navigate to service
cd services/api-gateway

# Service-specific commands
make build
make test
make run
make docker-build
```

## Docker Integration

### Service Dockerfiles
Each service has its own Dockerfile that:
- Copies the shared code
- Builds the service binary
- Creates a minimal runtime image

### Docker Compose
The `docker-compose.yml` orchestrates all services:
- Infrastructure services (PostgreSQL, Kafka)
- Application services
- Monitoring services

## Development Workflow

### 1. Working on a Single Service
```bash
# Navigate to service
cd services/api-gateway

# Install dependencies
make deps

# Run tests
make test

# Build and run
make run
```

### 2. Working on Shared Code
```bash
# Navigate to shared
cd shared

# Make changes
# Update tests

# Update services that use shared code
cd ../services/api-gateway
go mod tidy
```

### 3. Adding a New Service
1. Create new directory in `services/`
2. Copy structure from existing service
3. Update `go.mod` with required dependencies
4. Add to top-level `Makefile`
5. Update `docker-compose.yml`

## Migration Guide

### From Single-Project Structure

1. **Run the migration script**:
   ```bash
   ./scripts/migrate-to-monorepo.sh
   ```

2. **Review the new structure**:
   - Check that all code was moved correctly
   - Verify import paths were updated
   - Test building individual services

3. **Update your workflow**:
   - Use service-specific commands for development
   - Use top-level commands for building all services
   - Update CI/CD pipelines if needed

### Manual Migration Steps

If you prefer to migrate manually:

1. Create the new directory structure
2. Move code from `internal/` to `services/`
3. Move shared code to `shared/`
4. Create `go.mod` files for each service
5. Update import paths in Go files
6. Create service-specific `main.go` files
7. Update Dockerfiles and Makefiles

## Best Practices

### 1. Shared Code Guidelines
- Keep shared code minimal and focused
- Avoid service-specific logic in shared packages
- Use interfaces for shared contracts
- Version shared code carefully

### 2. Service Development
- Each service should be independently deployable
- Use dependency injection for external dependencies
- Write comprehensive tests for each service
- Follow consistent naming conventions

### 3. Dependency Management
- Keep service dependencies minimal
- Use the shared package for common functionality
- Regularly update dependencies
- Use `go mod tidy` to clean up unused dependencies

### 4. Testing Strategy
- Unit tests within each service
- Integration tests for service boundaries
- End-to-end tests for complete workflows
- Shared test utilities in the shared package

## Troubleshooting

### Common Issues

1. **Import Path Errors**
   - Ensure `go.mod` files are correct
   - Check that `replace` directives are in place
   - Run `go mod tidy` in each service

2. **Build Failures**
   - Check that shared code is accessible
   - Verify Docker build context includes shared code
   - Ensure all dependencies are declared

3. **Test Failures**
   - Run tests in each service directory
   - Check that shared test utilities are available
   - Verify mock implementations are correct

### Getting Help

- Check the service-specific README files
- Review the migration script for reference
- Look at existing services as examples
- Check the backup directory for the old structure

## Future Considerations

### Splitting to Polyrepos
If the project grows significantly, consider splitting to separate repositories:

1. **Shared Library Repository**
   - Extract shared code to its own repository
   - Publish as a Go module
   - Version and release independently

2. **Service Repositories**
   - Each service in its own repository
   - Independent CI/CD pipelines
   - Service-specific documentation

3. **Infrastructure Repository**
   - Docker Compose configurations
   - Kubernetes manifests
   - Deployment scripts

### Migration Path
The multi-project structure makes it easy to split later:
1. Extract shared code to separate repository
2. Update service dependencies to use published module
3. Move each service to its own repository
4. Update CI/CD and deployment configurations 