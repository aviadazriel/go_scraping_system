#!/bin/bash

# Go Scraping Project - Example Scripts
# This script contains example API calls and common operations

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL for API
API_BASE="http://localhost:8082"

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

# Function to check if service is running
check_service() {
    if curl -s -f "$API_BASE/health" > /dev/null; then
        print_success "API Gateway is running"
        return 0
    else
        print_error "API Gateway is not running. Please start the services first."
        return 1
    fi
}

# Function to create a URL
create_url() {
    local url="$1"
    local frequency="${2:-hourly}"
    local description="${3:-Test URL}"
    
    print_status "Creating URL: $url"
    
    response=$(curl -s -X POST "$API_BASE/api/v1/urls" \
        -H "Content-Type: application/json" \
        -d "{
            \"url\": \"$url\",
            \"frequency\": \"$frequency\",
            \"description\": \"$description\"
        }")
    
    if echo "$response" | grep -q "id"; then
        url_id=$(echo "$response" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
        print_success "URL created with ID: $url_id"
        echo "$url_id"
    else
        print_error "Failed to create URL: $response"
        return 1
    fi
}

# Function to list all URLs
list_urls() {
    print_status "Listing all URLs..."
    curl -s "$API_BASE/api/v1/urls" | jq '.' 2>/dev/null || curl -s "$API_BASE/api/v1/urls"
}

# Function to trigger scraping for a URL
trigger_url() {
    local url_id="$1"
    print_status "Triggering scraping for URL ID: $url_id"
    curl -s -X POST "$API_BASE/api/v1/urls/$url_id/trigger"
    print_success "Scraping triggered for URL ID: $url_id"
}

# Function to check Kafka topics
check_kafka() {
    print_status "Checking Kafka topics..."
    docker exec scraping_kafka kafka-topics --bootstrap-server localhost:9092 --list 2>/dev/null || {
        print_warning "Kafka not accessible. Make sure services are running."
        return 1
    }
}

# Function to monitor Kafka messages
monitor_kafka() {
    local topic="${1:-scraping-tasks}"
    print_status "Monitoring Kafka topic: $topic"
    print_warning "Press Ctrl+C to stop monitoring"
    docker exec scraping_kafka kafka-console-consumer \
        --bootstrap-server localhost:9092 \
        --topic "$topic" \
        --from-beginning
}

# Function to check database
check_database() {
    print_status "Checking database..."
    psql -h localhost -U scraper -d scraper -c "SELECT COUNT(*) as url_count FROM urls;" 2>/dev/null || {
        print_error "Database not accessible. Check PostgreSQL connection."
        return 1
    }
}

# Main menu
show_menu() {
    echo
    echo "=== Go Scraping Project - Example Scripts ==="
    echo "1. Check service health"
    echo "2. Create a test URL"
    echo "3. List all URLs"
    echo "4. Trigger scraping for a URL"
    echo "5. Check Kafka topics"
    echo "6. Monitor Kafka messages"
    echo "7. Check database"
    echo "8. Run full example workflow"
    echo "9. Exit"
    echo
    read -p "Select an option (1-9): " choice
}

# Full example workflow
run_example_workflow() {
    print_status "Running full example workflow..."
    
    # Check service health
    if ! check_service; then
        return 1
    fi
    
    # Create test URLs
    print_status "Creating test URLs..."
    url1_id=$(create_url "https://news.ycombinator.com" "hourly" "Hacker News")
    url2_id=$(create_url "https://github.com/trending" "daily" "GitHub Trending")
    
    # List URLs
    echo
    list_urls
    
    # Trigger scraping
    echo
    trigger_url "$url1_id"
    trigger_url "$url2_id"
    
    # Check Kafka
    echo
    check_kafka
    
    # Check database
    echo
    check_database
    
    print_success "Example workflow completed!"
}

# Main script logic
main() {
    case "$1" in
        "health")
            check_service
            ;;
        "create")
            create_url "${2:-https://example.com}" "${3:-hourly}" "${4:-Test URL}"
            ;;
        "list")
            list_urls
            ;;
        "trigger")
            if [ -z "$2" ]; then
                print_error "URL ID required. Usage: $0 trigger <url_id>"
                exit 1
            fi
            trigger_url "$2"
            ;;
        "kafka")
            check_kafka
            ;;
        "monitor")
            monitor_kafka "$2"
            ;;
        "db")
            check_database
            ;;
        "workflow")
            run_example_workflow
            ;;
        *)
            # Interactive mode
            while true; do
                show_menu
                case $choice in
                    1)
                        check_service
                        ;;
                    2)
                        read -p "Enter URL: " url
                        read -p "Enter frequency (hourly/daily/weekly): " freq
                        read -p "Enter description: " desc
                        create_url "$url" "$freq" "$desc"
                        ;;
                    3)
                        list_urls
                        ;;
                    4)
                        read -p "Enter URL ID: " url_id
                        trigger_url "$url_id"
                        ;;
                    5)
                        check_kafka
                        ;;
                    6)
                        read -p "Enter topic name (default: scraping-tasks): " topic
                        monitor_kafka "${topic:-scraping-tasks}"
                        ;;
                    7)
                        check_database
                        ;;
                    8)
                        run_example_workflow
                        ;;
                    9)
                        print_status "Goodbye!"
                        exit 0
                        ;;
                    *)
                        print_error "Invalid option. Please try again."
                        ;;
                esac
                echo
                read -p "Press Enter to continue..."
            done
            ;;
    esac
}

# Check if jq is installed for JSON formatting
if ! command -v jq &> /dev/null; then
    print_warning "jq not found. Install it for better JSON formatting: brew install jq (macOS) or apt-get install jq (Ubuntu)"
fi

# Run main function with all arguments
main "$@" 