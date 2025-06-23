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
.PHONY: test
test: ## Run all tests
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v -short ./...

.PHONY: test-integration
test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	go test -v -run Integration ./...

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