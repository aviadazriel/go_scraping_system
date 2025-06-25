# URL Manager Service Makefile

.PHONY: build test clean run docker-build docker-run

# Service name
SERVICE_NAME := url-manager

# Build variables
BINARY_NAME := $(SERVICE_NAME)
BUILD_DIR := build
DOCKER_IMAGE := $(SERVICE_NAME):latest

# Go variables
GO := go
GOFMT := gofmt
GOLINT := golangci-lint

# Build the service
build:
	@echo "Building $(SERVICE_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Run the service
run: build
	@echo "Running $(SERVICE_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage for $(SERVICE_NAME)..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Lint code
lint:
	@echo "Linting code..."
	$(GOLINT) run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Docker build
docker-build:
	@echo "Building Docker image for $(SERVICE_NAME)..."
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run: docker-build
	@echo "Running $(SERVICE_NAME) in Docker..."
	docker run -p 8081:8081 $(DOCKER_IMAGE)

# Development setup
dev-setup: deps fmt lint test

# Production build
prod-build: clean deps test build 