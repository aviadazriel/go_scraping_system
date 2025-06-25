# Multi-Project Monorepo Makefile

.PHONY: help build test clean deps fmt lint migrate-up migrate-down migrate-status sqlc-generate sqlc-check create-kafka-topics

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
	@echo "Building and running API Gateway service..."
	@cd services/api-gateway && make run

# URL Manager service
url-manager:
	@echo "Building and running URL Manager service..."
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
	docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --create --if-not-exists --topic scraping-tasks --partitions 1 --replication-factor 1
	# Add more topics as needed (e.g. scraping-results, url-updates)