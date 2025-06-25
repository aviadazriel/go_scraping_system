# Multi-Project Monorepo Makefile

.PHONY: help build test clean deps fmt lint migrate-up migrate-down migrate-status sqlc-generate sqlc-check create-kafka-topics infrastructure-up infrastructure-down dev-local

# Default target
help:
	@echo "Available targets:"
	@echo "  build            - Build all services"
	@echo "  test             - Run tests for all services"
	@echo "  clean            - Clean all build artifacts"
	@echo "  deps             - Install dependencies for all services"
	@echo "  fmt              - Format code for all services"
	@echo "  lint             - Lint code for all services"
	@echo "  docker-build-all - Build Docker images for all services"
	@echo "  docker-run-all   - Run all services in Docker"
	@echo "  dev-setup-all    - Setup development environment for all services"
	@echo "  api-gateway      - Build and run API Gateway service"
	@echo "  url-manager      - Build and run URL Manager service"
	@echo "  infrastructure-up   - Start only infrastructure services (Kafka, Zookeeper, PostgreSQL)"
	@echo "  infrastructure-down - Stop infrastructure services"
	@echo "  dev-local        - Setup for local development (infrastructure + local services)"
	@echo ""
	@echo "Database commands:"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback last migration"
	@echo "  migrate-status   - Show migration status"
	@echo "  sqlc-generate    - Generate Go code from SQL"
	@echo "  sqlc-check       - Check SQLC configuration"
	@echo "  create-kafka-topics - Create required Kafka topics for deployment"

# Build all services
build: build-all

# Test all services
test: test-all

# Clean all build artifacts
clean: clean-all

# Install dependencies for all services
deps: deps-all

# Format code for all services
fmt: fmt-all

# Lint code for all services
lint: lint-all

# Database migration shortcuts
migrate-up: db-migrate-up
migrate-down: db-migrate-down
migrate-status: db-migrate-status

# Infrastructure only (for local development)
infrastructure-up:
	@echo "Starting infrastructure services (Kafka, Zookeeper, PostgreSQL)..."
	@docker-compose -f docker-compose.infrastructure.yml up -d
	@echo "Waiting for services to be ready..."
	@sleep 15
	@echo "Waiting for Kafka to be ready..."
	@until docker exec scraper-kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; do \
		echo "Waiting for Kafka..."; \
		sleep 5; \
	done
	@echo "Creating Kafka topics..."
	@make create-kafka-topics
	@echo "Running database migrations..."
	@make migrate-up
	@echo "Infrastructure is ready! You can now run services locally:"
	@echo "  make api-gateway    # Run API Gateway locally (port 8080)"
	@echo "  make url-manager    # Run URL Manager locally"
	@echo "  Kafka UI available at: http://localhost:8081"

infrastructure-down:
	@echo "Stopping infrastructure services..."
	@docker-compose -f docker-compose.infrastructure.yml down

# Local development setup
dev-local: infrastructure-up
	@echo "Setting up local development environment..."
	@make deps-all
	@echo "Local development environment is ready!"
	@echo "Infrastructure services are running in Docker."
	@echo "You can now run services locally:"
	@echo "  make api-gateway    # Run API Gateway locally"
	@echo "  make url-manager    # Run URL Manager locally"

# Build all services
build-all:
	@echo "Building all services..."
	@cd services/api-gateway && make build
	@cd services/url-manager && make build

# Test all services
test-all:
	@echo "Testing all services..."
	@cd services/api-gateway && make test
	@cd services/url-manager && make test

# Clean all build artifacts
clean-all:
	@echo "Cleaning all build artifacts..."
	@cd services/api-gateway && make clean
	@cd services/url-manager && make clean

# Build Docker images for all services
docker-build-all:
	@echo "Building Docker images for all services..."
	@cd services/api-gateway && make docker-build
	@cd services/url-manager && make docker-build

# Run all services in Docker
docker-run-all: docker-build-all
	@echo "Running all services in Docker..."
	@docker-compose up -d

# Setup development environment for all services
dev-setup-all:
	@echo "Setting up development environment for all services..."
	@cd shared && go mod download
	@cd services/api-gateway && make dev-setup
	@cd services/url-manager && make dev-setup

# API Gateway service
api-gateway:
	@echo "Building and running API Gateway service locally..."
	@cd services/api-gateway && make run

# URL Manager service
url-manager:
	@echo "Building and running URL Manager service locally..."
	@cd services/url-manager && make run

# Install dependencies for all services
deps-all:
	@echo "Installing dependencies for all services..."
	@cd shared && go mod download
	@cd services/api-gateway && make deps
	@cd services/url-manager && make deps

# Format code for all services
fmt-all:
	@echo "Formatting code for all services..."
	@cd shared && go fmt ./...
	@cd services/api-gateway && make fmt
	@cd services/url-manager && make fmt

# Lint code for all services
lint-all:
	@echo "Linting code for all services..."
	@cd services/api-gateway && make lint
	@cd services/url-manager && make lint

# Production build for all services
prod-build-all:
	@echo "Building all services for production..."
	@cd services/api-gateway && make prod-build
	@cd services/url-manager && make prod-build

# Database operations
db-setup:
	@echo "Setting up database..."
	@cd shared/database && make db-setup

db-migrate-up:
	@echo "Running database migrations..."
	@cd shared/database && make migrate-up

db-migrate-down:
	@echo "Rolling back last migration..."
	@cd shared/database && make migrate-down

db-migrate-status:
	@echo "Migration status:"
	@cd shared/database && make migrate-status

sqlc-generate:
	@echo "Generating SQLC code..."
	@cd shared/database && make sqlc-generate

sqlc-check:
	@echo "Checking SQLC configuration..."
	@cd shared/database && make sqlc-check

create-kafka-topics:
	docker exec scraper-kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic scraping-requests --partitions 1 --replication-factor 1
	docker exec scraper-kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic scraping-results --partitions 1 --replication-factor 1
	docker exec scraper-kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic url-updates --partitions 1 --replication-factor 1