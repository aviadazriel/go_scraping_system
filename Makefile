# Go Scraping Project Makefile

# Variables
BINARY_DIR=bin
DOCKER_REGISTRY=localhost:5000
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Services
SERVICES=url-manager scraper parser data-storage api-gateway

# Default target
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
.PHONY: dev-setup
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	mkdir -p $(BINARY_DIR)
	mkdir -p data/html
	mkdir -p configs

.PHONY: build
build: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		go build $(LDFLAGS) -o $(BINARY_DIR)/$$service ./cmd/$$service; \
	done

.PHONY: build-service
build-service: ## Build a specific service (usage: make build-service SERVICE=url-manager)
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE variable is required. Usage: make build-service SERVICE=url-manager"; \
		exit 1; \
	fi
	@echo "Building $(SERVICE)..."
	@go build $(LDFLAGS) -o $(BINARY_DIR)/$(SERVICE) ./cmd/$(SERVICE)

.PHONY: run-service
run-service: ## Run a specific service (usage: make run-service SERVICE=url-manager)
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE variable is required. Usage: make run-service SERVICE=url-manager"; \
		exit 1; \
	fi
	@echo "Running $(SERVICE)..."
	@go run ./cmd/$(SERVICE)

.PHONY: run-all
run-all: ## Run all services in separate terminals
	@echo "Starting all services..."
	@for service in $(SERVICES); do \
		echo "Starting $$service..."; \
		gnome-terminal --title="$$service" -- go run ./cmd/$$service & \
	done

# Infrastructure
.PHONY: infra-up
infra-up: ## Start infrastructure services (PostgreSQL, Kafka, etc.)
	@echo "Starting infrastructure services..."
	docker-compose up -d postgres zookeeper kafka

.PHONY: infra-down
infra-down: ## Stop infrastructure services
	@echo "Stopping infrastructure services..."
	docker-compose down

.PHONY: infra-logs
infra-logs: ## Show infrastructure logs
	docker-compose logs -f

.PHONY: infra-clean
infra-clean: ## Clean up infrastructure (removes volumes)
	@echo "Cleaning up infrastructure..."
	docker-compose down -v
	docker system prune -f

# Database
.PHONY: migrate
migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run ./cmd/migrate

.PHONY: db-reset
db-reset: ## Reset database (drops and recreates)
	@echo "Resetting database..."
	docker-compose down postgres
	docker volume rm go_scraping_project_postgres_data
	docker-compose up -d postgres
	@sleep 5
	@make migrate

# Testing
.PHONY: test test-unit test-integration test-all test-coverage test-race test-benchmark

# Run all tests
test: test-unit test-integration

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	go test -v ./tests/unit/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	go test -v ./tests/integration/...

# Run all tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./tests/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./tests/...

# Run benchmarks
test-benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./tests/...

# Run tests in parallel
test-parallel:
	@echo "Running tests in parallel..."
	go test -v -parallel 4 ./tests/...

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -count=1 ./tests/...

# Run tests and generate coverage badge
test-coverage-badge:
	@echo "Running tests and generating coverage badge..."
	go test -v -coverprofile=coverage.out ./tests/...
	@if command -v gocover-cobertura >/dev/null 2>&1; then \
		gocover-cobertura < coverage.out > coverage.xml; \
		echo "Coverage XML generated: coverage.xml"; \
	else \
		echo "gocover-cobertura not found. Install with: go install github.com/boumenot/gocover-cobertura@latest"; \
	fi

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html coverage.xml
	rm -rf testdata/

# Setup test environment
test-setup:
	@echo "Setting up test environment..."
	@if [ ! -d "testdata" ]; then mkdir -p testdata; fi
	@echo "Test environment ready"

# Run tests with specific tags
test-short:
	@echo "Running short tests..."
	go test -v -short ./tests/...

test-long:
	@echo "Running long tests..."
	go test -v -run "Test.*Long" ./tests/...

# Run tests for specific packages
test-url-manager:
	@echo "Running URL Manager tests..."
	go test -v ./internal/url-manager/... ./tests/unit/url_*.go

test-api-gateway:
	@echo "Running API Gateway tests..."
	go test -v ./internal/api-gateway/... ./tests/unit/api_*.go

test-database:
	@echo "Running database tests..."
	go test -v ./internal/database/... ./tests/unit/database_*.go

# Run tests with different build tags
test-race-detector:
	@echo "Running tests with race detector..."
	go test -v -race -tags=race ./tests/...

test-debug:
	@echo "Running tests with debug output..."
	go test -v -tags=debug ./tests/...

# Performance testing
test-performance:
	@echo "Running performance tests..."
	go test -v -bench=. -benchmem -run=^$$ ./tests/...

# Load testing (if you have load test files)
test-load:
	@echo "Running load tests..."
	@if [ -d "tests/load" ]; then \
		go test -v ./tests/load/...; \
	else \
		echo "No load tests found in tests/load/"; \
	fi

# Security testing
test-security:
	@echo "Running security tests..."
	@if [ -d "tests/security" ]; then \
		go test -v ./tests/security/...; \
	else \
		echo "No security tests found in tests/security/"; \
	fi

# Test database setup
test-db-setup:
	@echo "Setting up test database..."
	@if command -v docker >/dev/null 2>&1; then \
		docker run --rm -d --name scraper-test-db \
			-e POSTGRES_DB=scraper_test \
			-e POSTGRES_USER=scraper \
			-e POSTGRES_PASSWORD=scraper_password \
			-p 5433:5432 \
			postgres:15; \
		echo "Test database started on port 5433"; \
		echo "Waiting for database to be ready..."; \
		sleep 10; \
	else \
		echo "Docker not found. Please install Docker or use existing PostgreSQL instance."; \
	fi

test-db-cleanup:
	@echo "Cleaning up test database..."
	@if command -v docker >/dev/null 2>&1; then \
		docker stop scraper-test-db 2>/dev/null || true; \
		docker rm scraper-test-db 2>/dev/null || true; \
		echo "Test database cleaned up"; \
	else \
		echo "Docker not found. Please clean up manually."; \
	fi

# Test with different Go versions (if you have multiple Go versions)
test-go-versions:
	@echo "Testing with different Go versions..."
	@for version in 1.21 1.22 1.23; do \
		echo "Testing with Go $$version..."; \
		if command -v go$$version >/dev/null 2>&1; then \
			go$$version test -v ./tests/unit/...; \
		else \
			echo "Go $$version not found, skipping..."; \
		fi; \
	done

# Test CI/CD pipeline simulation
test-ci:
	@echo "Running CI/CD pipeline tests..."
	@echo "1. Running linter..."
	@make lint
	@echo "2. Running unit tests..."
	@make test-unit
	@echo "3. Running integration tests..."
	@make test-integration
	@echo "4. Running security tests..."
	@make test-security
	@echo "5. Building application..."
	@make build
	@echo "CI/CD pipeline tests completed successfully!"

# Test documentation
test-docs:
	@echo "Testing documentation..."
	@if [ -f "README.md" ]; then \
		echo "README.md exists"; \
	else \
		echo "Warning: README.md not found"; \
	fi
	@if [ -d "docs" ]; then \
		echo "Documentation directory exists"; \
		ls -la docs/; \
	else \
		echo "Warning: docs/ directory not found"; \
	fi

# Test configuration
test-config:
	@echo "Testing configuration files..."
	@if [ -f "configs/config.yaml" ]; then \
		echo "config.yaml exists and is valid"; \
		yq eval '.' configs/config.yaml >/dev/null 2>&1 || echo "Warning: config.yaml is not valid YAML"; \
	else \
		echo "Warning: configs/config.yaml not found"; \
	fi

# Test all components
test-full: test-setup test-db-setup test test-docs test-config test-db-cleanup
	@echo "Full test suite completed!"

# Code Quality
.PHONY: fmt
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	go mod tidy
	go mod verify

# Docker
.PHONY: docker-build
docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@for service in $(SERVICES); do \
		echo "Building $$service Docker image..."; \
		docker build -t $(DOCKER_REGISTRY)/$$service:$(VERSION) -f deployments/docker/$$service.Dockerfile .; \
	done

.PHONY: docker-push
docker-push: ## Push Docker images
	@echo "Pushing Docker images..."
	@for service in $(SERVICES); do \
		echo "Pushing $$service Docker image..."; \
		docker push $(DOCKER_REGISTRY)/$$service:$(VERSION); \
	done

.PHONY: docker-run
docker-run: ## Run services with Docker Compose
	@echo "Running services with Docker Compose..."
	docker-compose -f docker-compose.yml -f docker-compose.services.yml up -d

# Monitoring
.PHONY: monitoring-up
monitoring-up: ## Start monitoring stack (Prometheus, Grafana, Jaeger)
	@echo "Starting monitoring stack..."
	docker-compose up -d prometheus grafana jaeger

.PHONY: monitoring-down
monitoring-down: ## Stop monitoring stack
	@echo "Stopping monitoring stack..."
	docker-compose down prometheus grafana jaeger

# Cleanup
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: ## Clean everything (build artifacts, Docker images, volumes)
	@echo "Cleaning everything..."
	@make clean
	@make infra-clean
	docker system prune -af

# Development helpers
.PHONY: watch
watch: ## Watch for changes and rebuild (requires air)
	@echo "Watching for changes..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

.PHONY: install-tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-delve/delve/cmd/dlv@latest

.PHONY: generate
generate: ## Generate code (protobuf, mocks, etc.)
	@echo "Generating code..."
	@if [ -f "scripts/generate.sh" ]; then \
		./scripts/generate.sh; \
	else \
		echo "No generate script found"; \
	fi

# Health checks
.PHONY: health-check
health-check: ## Check health of all services
	@echo "Checking service health..."
	@curl -f http://localhost:8080/health || echo "API Gateway: DOWN"
	@curl -f http://localhost:8081/health || echo "URL Manager: DOWN"
	@curl -f http://localhost:8082/health || echo "Scraper: DOWN"
	@curl -f http://localhost:8083/health || echo "Parser: DOWN"
	@curl -f http://localhost:8084/health || echo "Data Storage: DOWN"

# Documentation
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/api-gateway/main.go; \
	else \
		echo "Swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Release
.PHONY: release
release: ## Create a release
	@echo "Creating release..."
	@if [ -z "$(TAG)" ]; then \
		echo "Error: TAG variable is required. Usage: make release TAG=v1.0.0"; \
		exit 1; \
	fi
	git tag $(TAG)
	git push origin $(TAG)
	@make docker-build
	@make docker-push 


.PHONY: migrate-up
migrate-up:
	goose -dir ./sql/schema postgres "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir ./sql/schema postgres "$(DATABASE_URL)" down

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate