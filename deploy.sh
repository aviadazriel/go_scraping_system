#!/bin/bash

# URL Manager Service Deployment Script
# This script deploys the URL Manager service with Kafka and connects to your local PostgreSQL

set -e

echo "🚀 Starting URL Manager Service Deployment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker version > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please install it and try again."
    exit 1
fi

# Check if PostgreSQL is running locally
print_status "Checking local PostgreSQL connection..."
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    print_warning "PostgreSQL is not running on localhost:5432"
    print_warning "Please make sure PostgreSQL is running and accessible"
    print_warning "You can start it with: brew services start postgresql (macOS) or sudo systemctl start postgresql (Linux)"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    print_success "PostgreSQL is running on localhost:5432"
fi

# Check if the scraper database exists
print_status "Checking if 'scraper' database exists..."
if ! psql -h localhost -U scraper -d scraper -c "SELECT 1;" > /dev/null 2>&1; then
    print_warning "Database 'scraper' or user 'scraper' does not exist"
    print_status "Creating database and user..."
    
    # Create user and database (requires superuser privileges)
    psql -h localhost -U postgres -c "CREATE USER scraper WITH PASSWORD 'scraper_password';" || true
    psql -h localhost -U postgres -c "CREATE DATABASE scraper OWNER scraper;" || true
    psql -h localhost -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE scraper TO scraper;" || true
    
    print_success "Database and user created"
else
    print_success "Database 'scraper' exists and is accessible"
fi

# Run database migrations
print_status "Running database migrations..."
if command -v goose &> /dev/null; then
    goose -dir sql/migrations postgres "host=localhost port=5432 user=scraper password=scraper_password dbname=scraper sslmode=disable" up
    print_success "Database migrations completed"
else
    print_warning "Goose migration tool not found. Skipping migrations."
    print_warning "Please run migrations manually: make migrate-up"
fi

# Build and start services
print_status "Building and starting services..."

# Stop any existing containers
docker-compose -f docker-compose.local.yml down --remove-orphans

# Build and start services
docker-compose -f docker-compose.local.yml up -d --build

# Wait for services to be healthy
print_status "Waiting for services to be healthy..."

# Wait for Kafka
print_status "Waiting for Kafka to be ready..."
timeout=120
counter=0
while [ $counter -lt $timeout ]; do
    if docker-compose -f docker-compose.local.yml exec -T kafka kafka-topics --bootstrap-server localhost:9092 --list > /dev/null 2>&1; then
        print_success "Kafka is ready"
        break
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done

if [ $counter -eq $timeout ]; then
    print_error "Kafka failed to start within $timeout seconds"
    docker-compose -f docker-compose.local.yml logs kafka
    exit 1
fi

# Wait for URL Manager
print_status "Waiting for URL Manager to be ready..."
timeout=60
counter=0
while [ $counter -lt $timeout ]; do
    if curl -f http://localhost:8081/health > /dev/null 2>&1; then
        print_success "URL Manager is ready"
        break
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done

if [ $counter -eq $timeout ]; then
    print_error "URL Manager failed to start within $timeout seconds"
    docker-compose -f docker-compose.local.yml logs url-manager
    exit 1
fi

# Wait for API Gateway
print_status "Waiting for API Gateway to be ready..."
timeout=60
counter=0
while [ $counter -lt $timeout ]; then
    if curl -f http://localhost:8082/health > /dev/null 2>&1; then
        print_success "API Gateway is ready"
        break
    fi
    sleep 2
    counter=$((counter + 2))
    echo -n "."
done

if [ $counter -eq $timeout ]; then
    print_error "API Gateway failed to start within $timeout seconds"
    docker-compose -f docker-compose.local.yml logs api-gateway
    exit 1
fi

# Display service information
echo
print_success "🎉 Deployment completed successfully!"
echo
echo "📋 Service Information:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 API Gateway:     http://localhost:8082"
echo "🔧 URL Manager:     http://localhost:8081"
echo "📊 Kafka UI:        http://localhost:8080"
echo "🗄️  PostgreSQL:      localhost:5432"
echo "📨 Kafka:           localhost:9092"
echo
echo "🔍 Health Checks:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "API Gateway:     curl http://localhost:8082/health"
echo "URL Manager:     curl http://localhost:8081/health"
echo
echo "📝 Useful Commands:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "View logs:        docker-compose -f docker-compose.local.yml logs -f"
echo "Stop services:    docker-compose -f docker-compose.local.yml down"
echo "Restart services: docker-compose -f docker-compose.local.yml restart"
echo
echo "🧪 Test the API:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Create URL:       curl -X POST http://localhost:8082/api/v1/urls \\"
echo "                   -H 'Content-Type: application/json' \\"
echo "                   -d '{\"url\":\"https://example.com\",\"frequency\":\"1h\"}'"
echo
print_success "Deployment script completed!" 