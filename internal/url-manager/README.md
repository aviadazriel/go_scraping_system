# URL Manager Service

## Overview

The URL Manager is a **background microservice** responsible for scheduling and distributing web scraping tasks. It operates as a daemon process that continuously monitors the database for URLs that are due for scraping and sends them to Kafka topics for processing by other services.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Database      │    │  URL Manager    │    │     Kafka       │
│   (PostgreSQL)  │◄──►│  (Background    │───►│   (Topics)      │
│                 │    │   Service)      │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Responsibilities

### 1. **URL Scheduling**
- Periodically scans the database for URLs due for scraping
- Calculates next scrape times based on frequency settings
- Updates URL status and scheduling information

### 2. **Task Distribution**
- Creates scraping tasks for due URLs
- Sends tasks to Kafka topics for processing
- Ensures reliable message delivery

### 3. **Status Management**
- Updates URL status (pending → in_progress → completed/failed)
- Tracks retry counts and last scraped times
- Manages scheduling metadata

## Service Components

### Core Services

#### `URLSchedulerService`
- **Purpose**: Main scheduling engine that runs continuously
- **Functionality**:
  - Runs every 30 seconds to check for due URLs
  - Processes URLs scheduled for scraping within a time window
  - Creates and sends Kafka messages for each task
  - Updates database with new scheduling information

#### `URLRepository`
- **Purpose**: Data access layer for URL operations
- **Functionality**:
  - Retrieves URLs scheduled for scraping
  - Updates URL status and timing information
  - Manages retry counts and metadata

### Data Models

#### `Frequency` Types
```go
const (
    Frequency30Seconds = "30s"
    Frequency1Minute   = "1m"
    Frequency5Minutes  = "5m"
    Frequency15Minutes = "15m"
    Frequency30Minutes = "30m"
    Frequency1Hour     = "1h"
    Frequency6Hours    = "6h"
    Frequency12Hours   = "12h"
    Frequency1Day      = "1d"
    Frequency1Week     = "1w"
)
```

## Operation Flow

### 1. **Scheduling Loop**
```
Every 30 seconds:
├── Query database for URLs due in next 5 minutes
├── For each due URL:
│   ├── Create scraping task
│   ├── Send to Kafka topic
│   ├── Update URL status to "in_progress"
│   ├── Calculate next scrape time
│   └── Update database
└── Log processing results
```

### 2. **Task Creation**
```go
task := &domain.ScrapingTask{
    ID:        uuid.New(),
    URLID:     url.ID,
    URL:       url.Url,
    Status:    domain.URLStatusPending,
    Attempt:   1,
    CreatedAt: time.Now(),
}
```

### 3. **Kafka Message**
```go
message := domain.NewScrapingTaskMessage(task, uuid.New().String())
producer.SendMessage(ctx, domain.TopicScrapingTasks, message)
```

## Configuration

### Environment Variables
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=scraping_db
DB_USER=scraping_user
DB_PASSWORD=scraping_password

# Kafka
KAFKA_BROKERS=kafka:29092

# Logging
LOG_LEVEL=info
```

### Scheduling Parameters
- **Check Interval**: 30 seconds (configurable)
- **Time Window**: ±5 minutes around current time
- **Batch Size**: Up to 100 URLs per cycle
- **Retry Logic**: Built into URL frequency calculation

## Database Schema

### URLs Table
```sql
CREATE TABLE urls (
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
    frequency TEXT NOT NULL,
    status TEXT NOT NULL,
    next_scrape_at TIMESTAMP,
    last_scraped_at TIMESTAMP,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    -- ... other fields
);
```

### Key Queries
- `GetURLsScheduledForScraping`: Find URLs due for scraping
- `UpdateURLStatus`: Update URL processing status
- `UpdateNextScrapeTime`: Schedule next scrape
- `IncrementRetryCount`: Track retry attempts

## Kafka Topics

### `scraping-tasks`
- **Purpose**: Distribute scraping tasks to scraper services
- **Message Format**: `ScrapingTask` with URL and metadata
- **Consumers**: Scraper services

### Message Structure
```json
{
  "id": "task-uuid",
  "url_id": "url-uuid", 
  "url": "https://example.com",
  "status": "pending",
  "attempt": 1,
  "created_at": "2024-01-01T00:00:00Z"
}
```

## Monitoring & Observability

### Logging
- **Structured Logging**: JSON format with correlation IDs
- **Log Levels**: INFO, WARN, ERROR
- **Key Events**:
  - Service startup/shutdown
  - URL processing cycles
  - Kafka message delivery
  - Database operations

### Metrics (Future)
- URLs processed per cycle
- Kafka message delivery success rate
- Database query performance
- Scheduling accuracy

### Health Checks
- Database connectivity
- Kafka producer status
- Service uptime

## Deployment

### Docker
```bash
# Build and run
docker-compose up url-manager

# Or build individually
docker build -f cmd/url-manager/Dockerfile -t url-manager .
docker run --env-file .env url-manager
```

### Kubernetes (Future)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: url-manager
spec:
  replicas: 1  # Single instance for scheduling
  selector:
    matchLabels:
      app: url-manager
  template:
    spec:
      containers:
      - name: url-manager
        image: url-manager:latest
        env:
        - name: DB_HOST
          value: postgres-service
        - name: KAFKA_BROKERS
          value: kafka-service:9092
```

## Error Handling

### Database Errors
- **Connection Issues**: Retry with exponential backoff
- **Query Failures**: Log error and continue with next URL
- **Transaction Failures**: Rollback and retry

### Kafka Errors
- **Producer Failures**: Log error and mark URL for retry
- **Message Delivery**: Use Kafka's built-in retry mechanism
- **Topic Issues**: Alert and stop processing

### URL Processing Errors
- **Invalid URLs**: Skip and log warning
- **Frequency Errors**: Use default frequency
- **Status Update Failures**: Continue processing

## Best Practices

### 1. **Idempotency**
- Use UUIDs for all operations
- Check URL status before processing
- Handle duplicate messages gracefully

### 2. **Resilience**
- Graceful shutdown handling
- Context cancellation support
- Resource cleanup on exit

### 3. **Performance**
- Batch database operations
- Efficient time-based queries
- Connection pooling

### 4. **Observability**
- Structured logging with correlation IDs
- Metrics for key operations
- Health check endpoints

## Troubleshooting

### Common Issues

#### 1. **URLs Not Being Processed**
- Check database connectivity
- Verify URL status and next_scrape_at values
- Review scheduler logs for errors

#### 2. **Kafka Message Failures**
- Verify Kafka broker connectivity
- Check topic configuration
- Review producer logs

#### 3. **High CPU/Memory Usage**
- Check scheduling frequency
- Review batch sizes
- Monitor database query performance

### Debug Commands
```bash
# Check service logs
docker logs scraping_url_manager

# Check database
docker exec -it scraping_postgres psql -U scraping_user -d scraping_db

# Check Kafka topics
docker exec -it scraping_kafka kafka-topics --list --bootstrap-server localhost:9092
```

## Future Enhancements

### 1. **Advanced Scheduling**
- Dynamic frequency adjustment
- Priority-based scheduling
- Timezone support

### 2. **Scalability**
- Horizontal scaling with leader election
- Distributed scheduling
- Load balancing

### 3. **Monitoring**
- Prometheus metrics
- Grafana dashboards
- Alerting rules

### 4. **Features**
- Manual trigger endpoints
- Bulk operations
- Scheduling analytics 