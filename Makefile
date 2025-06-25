# Multi-Project Monorepo Makefile

.PHONY: help build-all test-all clean-all docker-build-all docker-run-all dev-setup-all

# Default target
help:
	@echo "Available targets:"
	@echo "  build-all        - Build all services"
	@echo "  test-all         - Run tests for all services"
	@echo "  clean-all        - Clean all build artifacts"
	@echo "  docker-build-all - Build Docker images for all services"
	@echo "  docker-run-all   - Run all services in Docker"
	@echo "  dev-setup-all    - Setup development environment for all services"
	@echo "  api-gateway      - Build and run API Gateway service"
	@echo "  url-manager      - Build and run URL Manager service"
	@echo ""
	@echo "Database commands:"
	@echo "  db-setup         - Setup database (migrate + generate)"
	@echo "  db-migrate-up    - Run database migrations"
	@echo "  db-migrate-down  - Rollback last migration"
	@echo "  db-migrate-status- Show migration status"
	@echo "  sqlc-generate    - Generate Go code from SQL"
	@echo "  sqlc-check       - Check SQLC configuration"

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