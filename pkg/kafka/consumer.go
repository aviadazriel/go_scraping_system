package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/internal/domain"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// MessageHandler is a function type for handling Kafka messages
type MessageHandler func(ctx context.Context, message *domain.KafkaMessage) error

// Consumer represents a Kafka consumer using kafka-go
type Consumer struct {
	readers  map[string]*kafka.Reader
	config   *config.Config
	logger   *logrus.Logger
	handlers map[domain.MessageType]MessageHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, log *logrus.Logger) (*Consumer, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		readers:  make(map[string]*kafka.Reader),
		config:   cfg,
		logger:   log,
		handlers: make(map[domain.MessageType]MessageHandler),
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// RegisterHandler registers a message handler for a specific message type
func (c *Consumer) RegisterHandler(messageType domain.MessageType, handler MessageHandler) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[messageType] = handler
}

// getHandler returns the handler for a message type
func (c *Consumer) getHandler(messageType domain.MessageType) (MessageHandler, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	handler, exists := c.handlers[messageType]
	return handler, exists
}

// Consume starts consuming messages from the specified topics
func (c *Consumer) Consume(topics []string) error {
	c.logger.WithField("topics", topics).Info("Starting Kafka consumer")

	var wg sync.WaitGroup

	for _, topic := range topics {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:         c.config.Kafka.Brokers,
			Topic:           topic,
			GroupID:         c.config.Kafka.GroupID,
			MinBytes:        10e3, // 10KB
			MaxBytes:        10e6, // 10MB
			MaxWait:         1 * time.Second,
			ReadLagInterval: -1,
			Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
				c.logger.Debugf(msg, args...)
			}),
		})

		c.mu.Lock()
		c.readers[topic] = reader
		c.mu.Unlock()

		wg.Add(1)
		go func(topic string, reader *kafka.Reader) {
			defer wg.Done()
			defer reader.Close()
			c.consumeTopic(topic, reader)
		}(topic, reader)
	}

	// Wait for context cancellation
	<-c.ctx.Done()

	// Wait for all goroutines to finish
	wg.Wait()

	return c.ctx.Err()
}

// consumeTopic consumes messages from a specific topic
func (c *Consumer) consumeTopic(topic string, reader *kafka.Reader) {
	c.logger.WithField("topic", topic).Info("Starting to consume topic")

	for {
		select {
		case <-c.ctx.Done():
			c.logger.WithField("topic", topic).Info("Stopping topic consumption")
			return
		default:
			msg, err := reader.ReadMessage(c.ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				c.logger.WithError(err).WithField("topic", topic).Error("Failed to read message")
				time.Sleep(100 * time.Millisecond) // Brief pause before retry
				continue
			}

			c.logger.WithFields(logrus.Fields{
				"topic":     topic,
				"partition": msg.Partition,
				"offset":    msg.Offset,
				"key":       string(msg.Key),
			}).Debug("Received message")

			// Parse the message
			var kafkaMessage domain.KafkaMessage
			if err := json.Unmarshal(msg.Value, &kafkaMessage); err != nil {
				c.logger.WithError(err).Error("Failed to unmarshal message")
				continue
			}

			// Process the message
			if err := c.processMessage(c.ctx, &kafkaMessage, &msg); err != nil {
				c.logger.WithError(err).Error("Failed to process message")
				// Don't mark as consumed, let it be retried
				continue
			}
		}
	}
}

// processMessage processes a single message
func (c *Consumer) processMessage(ctx context.Context, message *domain.KafkaMessage, kafkaMsg *kafka.Message) error {
	handler, exists := c.getHandler(message.Type)
	if !exists {
		return fmt.Errorf("no handler registered for message type: %s", message.Type)
	}

	// Add correlation ID to context if present
	if message.Metadata.CorrelationID != "" {
		ctx = context.WithValue(ctx, "correlation_id", message.Metadata.CorrelationID)
	}

	// Process the message with retry logic
	return c.processWithRetry(ctx, message, handler, kafkaMsg)
}

// processWithRetry processes a message with retry logic
func (c *Consumer) processWithRetry(ctx context.Context, message *domain.KafkaMessage, handler MessageHandler, kafkaMsg *kafka.Message) error {
	maxRetries := c.config.Kafka.RetryMaxAttempts

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := handler(ctx, message)
		if err == nil {
			return nil
		}

		c.logger.WithFields(logrus.Fields{
			"message_id":   message.ID,
			"message_type": message.Type,
			"attempt":      attempt + 1,
			"max_retries":  maxRetries,
			"error":        err.Error(),
		}).Warn("Message processing failed, retrying...")

		// If this is the last attempt, send to dead letter queue
		if attempt == maxRetries {
			return c.sendToDeadLetter(message, err, kafkaMsg)
		}

		// Wait before retrying
		backoff := c.config.Kafka.RetryBackoff * time.Duration(attempt+1)
		select {
		case <-time.After(backoff):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// sendToDeadLetter sends a failed message to the dead letter queue
func (c *Consumer) sendToDeadLetter(message *domain.KafkaMessage, err error, kafkaMsg *kafka.Message) error {
	// Parse message ID as UUID
	messageID, parseErr := uuid.Parse(message.ID)
	if parseErr != nil {
		c.logger.WithError(parseErr).Error("Failed to parse message ID as UUID")
		messageID = uuid.New() // Fallback to new UUID
	}

	deadLetterMsg := &domain.DeadLetterMessage{
		ID:         messageID,
		Topic:      kafkaMsg.Topic,
		Partition:  int32(kafkaMsg.Partition),
		Offset:     kafkaMsg.Offset,
		Key:        kafkaMsg.Key,
		Value:      kafkaMsg.Value,
		Error:      err.Error(),
		RetryCount: message.Metadata.RetryCount,
		MaxRetries: c.config.Kafka.RetryMaxAttempts,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	// Here you would typically save the dead letter message to the database
	// For now, we'll just log it
	c.logger.WithFields(logrus.Fields{
		"message_id":  deadLetterMsg.ID,
		"topic":       deadLetterMsg.Topic,
		"error":       deadLetterMsg.Error,
		"retry_count": deadLetterMsg.RetryCount,
		"max_retries": deadLetterMsg.MaxRetries,
	}).Error("Message sent to dead letter queue")

	return err
}

// Close closes the consumer and all readers
func (c *Consumer) Close() error {
	c.cancel()

	c.mu.Lock()
	defer c.mu.Unlock()

	var lastErr error
	for topic, reader := range c.readers {
		if err := reader.Close(); err != nil {
			c.logger.WithError(err).WithField("topic", topic).Error("Failed to close reader")
			lastErr = err
		}
	}

	return lastErr
}

// HealthCheck performs a consumer health check
func (c *Consumer) HealthCheck(ctx context.Context) error {
	// Check if consumer is still running
	select {
	case <-c.ctx.Done():
		return fmt.Errorf("consumer context is cancelled")
	default:
		return nil
	}
}
