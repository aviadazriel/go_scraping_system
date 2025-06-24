# Kafka Package

This package provides Kafka producer and consumer functionality for the scraping microservices system using **kafka-go**. It includes both low-level implementations and convenient wrappers that implement domain interfaces.

## Technology Choice

We use **kafka-go** instead of Sarama because:
- **Simpler, more Go-idiomatic API**
- **Lower memory usage** and better GC characteristics
- **Built for Go** from the ground up
- **Active development** and modern Go features
- **Context support** throughout the API

## Components

### Core Components

#### `Producer`
- **Purpose**: Low-level Kafka producer implementation using kafka-go
- **Features**:
  - **Topic-based writers** with connection pooling
  - **Synchronous message production** for reliability
  - **Automatic batching** for performance
  - **Structured logging** with correlation IDs
  - **Health checks** and error handling

#### `Consumer`
- **Purpose**: Low-level Kafka consumer implementation using kafka-go
- **Features**:
  - **Consumer group support** with automatic rebalancing
  - **Message retry logic** with exponential backoff
  - **Dead letter queue** handling
  - **Concurrent topic consumption** with goroutines
  - **Graceful shutdown** with context cancellation

### Wrapper Components

#### `KafkaProducerWrapper`
- **Purpose**: Wraps the `Producer` to implement `domain.KafkaProducer` interface
- **Benefits**:
  - Clean separation between domain interfaces and implementation
  - Consistent interface across all services
  - Easy testing with mocks
  - Type safety

#### `KafkaConsumerWrapper`
- **Purpose**: Wraps the `Consumer` to provide a simplified interface
- **Benefits**:
  - Cleaner API for service consumption
  - Consistent error handling
  - Simplified handler registration

## Usage

### Producer Usage

```go
// Create a producer
producer, err := kafka.NewProducer(cfg, log)
if err != nil {
    log.Fatalf("Failed to create producer: %v", err)
}
defer producer.Close()

// Wrap it for domain interface compatibility
kafkaProducer := kafka.NewKafkaProducerWrapper(producer)

// Use in services
schedulerService := services.NewURLSchedulerService(urlRepo, kafkaProducer, log)
```

### Consumer Usage

```go
// Create a consumer
consumer, err := kafka.NewConsumer(cfg, log)
if err != nil {
    log.Fatalf("Failed to create consumer: %v", err)
}
defer consumer.Close()

// Wrap it for simplified usage
kafkaConsumer := kafka.NewKafkaConsumerWrapper(consumer)

// Register handlers
kafkaConsumer.RegisterHandler(domain.MessageTypeScrapingTask, func(ctx context.Context, msg *domain.KafkaMessage) error {
    // Handle scraping task
    return nil
})

// Start consuming
err = kafkaConsumer.Consume([]string{domain.TopicScrapingTasks})
```

## Configuration

### Producer Configuration
```go
writer := kafka.NewWriter(kafka.WriterConfig{
    Brokers:      []string{"localhost:9092"},
    Topic:        "scraping-tasks",
    BatchSize:    100,                    // Messages per batch
    BatchTimeout: 10 * time.Millisecond,  // Max wait for batch
    RequiredAcks: -1,                     // Require all replicas
    Async:        false,                  // Synchronous for reliability
})
```

### Consumer Configuration
```go
reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers:         []string{"localhost:9092"},
    Topic:           "scraping-tasks",
    GroupID:         "scraper-group",
    MinBytes:        10e3,                // 10KB min message size
    MaxBytes:        10e6,                // 10MB max message size
    MaxWait:         1 * time.Second,     // Max wait for messages
    ReadLagInterval: -1,                  // Disable lag monitoring
})
```

## Message Types

### Scraping Tasks
- **Topic**: `scraping-tasks`
- **Purpose**: Distribute scraping tasks to scraper services
- **Producer**: URL Manager Service
- **Consumers**: Scraper Services

### Scraped Data
- **Topic**: `scraped-data`
- **Purpose**: Share scraped HTML content
- **Producer**: Scraper Services
- **Consumers**: Parser Services

### Parsed Data
- **Topic**: `parsed-data`
- **Purpose**: Share structured parsed data
- **Producer**: Parser Services
- **Consumers**: Data Storage Services

### Dead Letter Queue
- **Topic**: `dead-letter`
- **Purpose**: Handle failed messages
- **Producer**: All services
- **Consumers**: Dead Letter Handler

### Retry Queue
- **Topic**: `retry`
- **Purpose**: Retry failed messages
- **Producer**: All services
- **Consumers**: All services

## Error Handling

### Producer Errors
- **Connection Issues**: Automatic retry with exponential backoff
- **Message Serialization**: Log error and skip message
- **Topic Issues**: Alert and stop processing

### Consumer Errors
- **Message Processing**: Retry with backoff, then dead letter queue
- **Deserialization Errors**: Log and skip message
- **Handler Errors**: Retry up to max attempts

## Performance Optimizations

### Producer Optimizations
- **Connection pooling** per topic
- **Message batching** for throughput
- **Compression** for large messages
- **Async/sync modes** for different use cases

### Consumer Optimizations
- **Concurrent topic consumption** with goroutines
- **Batch message processing** where possible
- **Efficient memory usage** with streaming
- **Graceful shutdown** with context cancellation

## Best Practices

### 1. **Message Idempotency**
- Use UUIDs for message IDs
- Check message processing status before processing
- Handle duplicate messages gracefully

### 2. **Error Handling**
- Always check for errors after operations
- Use structured logging with correlation IDs
- Implement proper retry logic

### 3. **Resource Management**
- Always close producers and consumers
- Use defer statements for cleanup
- Monitor connection health

### 4. **Performance**
- Use appropriate batch sizes
- Configure compression for large messages
- Monitor producer/consumer lag

## Testing

### Mock Usage
```go
// Create a mock producer for testing
mockProducer := &MockKafkaProducer{}

// Use in tests
schedulerService := services.NewURLSchedulerService(urlRepo, mockProducer, log)
```

### Integration Testing
```go
// Use real Kafka for integration tests
producer, err := kafka.NewProducer(cfg, log)
consumer, err := kafka.NewConsumer(cfg, log)

// Test message flow
err = producer.SendScrapingTask(ctx, task)
// ... verify consumer receives message
```

## Monitoring

### Metrics
- Message production rate
- Message consumption rate
- Error rates
- Lag metrics
- Connection health

### Logging
- Structured JSON logging
- Correlation IDs for message tracing
- Error context and stack traces
- Performance metrics

## Migration from Sarama

This package was migrated from Sarama to kafka-go for:
- **Simpler API** and better Go integration
- **Lower memory usage** and better performance
- **Active development** and modern features
- **Better context support** throughout

### Key Differences
- **Simpler configuration** with fewer options
- **Topic-based writers** instead of single producer
- **Reader-based consumers** instead of consumer groups
- **Better error handling** with context cancellation

## Troubleshooting

### Common Issues

#### 1. **Producer Connection Failures**
- Check Kafka broker connectivity
- Verify broker configuration
- Check network connectivity

#### 2. **Consumer Group Issues**
- Verify group ID configuration
- Check partition assignment
- Monitor rebalancing events

#### 3. **Message Processing Failures**
- Check message format
- Verify handler implementation
- Review dead letter queue

### Debug Commands
```bash
# Check Kafka topics
kafka-topics --list --bootstrap-server localhost:9092

# Check consumer groups
kafka-consumer-groups --bootstrap-server localhost:9092 --list

# Check topic offsets
kafka-run-class kafka.tools.GetOffsetShell --bootstrap-server localhost:9092 --topic scraping-tasks
``` 