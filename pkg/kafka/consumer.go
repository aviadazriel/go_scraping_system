package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/internal/domain"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MessageHandler defines the interface for handling Kafka messages
type MessageHandler func(ctx context.Context, message *domain.KafkaMessage) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.ConsumerGroup
	config   *config.Config
	logger   *logrus.Logger
	handlers map[domain.MessageType]MessageHandler
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// ConsumerGroup represents a consumer group
type ConsumerGroup struct {
	consumer *Consumer
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, log *logrus.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = cfg.Kafka.EnableAutoCommit
	config.Consumer.Offsets.AutoCommit.Interval = cfg.Kafka.AutoCommitInterval
	config.Consumer.Offsets.CommitInterval = cfg.Kafka.AutoCommitInterval
	config.Consumer.Group.Session.Timeout = cfg.Kafka.SessionTimeout
	config.Consumer.Group.Heartbeat.Interval = cfg.Kafka.HeartbeatInterval
	config.Consumer.MaxProcessingTime = cfg.Kafka.MaxPollInterval
	config.Consumer.Fetch.Max = int32(cfg.Kafka.MaxPollRecords)
	config.Consumer.Return.Errors = true

	// Use the latest version
	config.Version = sarama.V2_8_0_0

	consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		consumer: consumer,
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

// Consume starts consuming messages from the specified topics
func (c *Consumer) Consume(topics []string) error {
	c.logger.WithField("topics", topics).Info("Starting Kafka consumer")

	for {
		select {
		case <-c.ctx.Done():
			return c.ctx.Err()
		default:
			handler := &ConsumerGroup{consumer: c}
			err := c.consumer.Consume(c.ctx, topics, handler)
			if err != nil {
				c.logger.WithError(err).Error("Error from consumer")
				return fmt.Errorf("consumer error: %w", err)
			}
		}
	}
}

// Setup is called when a new consumer group session is about to begin
func (cg *ConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	cg.consumer.logger.Info("Consumer group session setup")
	return nil
}

// Cleanup is called when a consumer group session has ended
func (cg *ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	cg.consumer.logger.Info("Consumer group session cleanup")
	return nil
}

// ConsumeClaim processes messages from a partition
func (cg *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				continue
			}

			cg.consumer.logger.WithFields(logrus.Fields{
				"topic":     message.Topic,
				"partition": message.Partition,
				"offset":    message.Offset,
				"key":       string(message.Key),
			}).Debug("Received message")

			// Parse the message
			var kafkaMessage domain.KafkaMessage
			if err := json.Unmarshal(message.Value, &kafkaMessage); err != nil {
				cg.consumer.logger.WithError(err).Error("Failed to unmarshal message")
				session.MarkMessage(message, "")
				continue
			}

			// Process the message
			if err := cg.processMessage(session.Context(), &kafkaMessage, message); err != nil {
				cg.consumer.logger.WithError(err).Error("Failed to process message")
				// Don't mark as consumed, let it be retried
				continue
			}

			// Mark message as consumed
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// processMessage processes a single message
func (cg *ConsumerGroup) processMessage(ctx context.Context, message *domain.KafkaMessage, saramaMessage *sarama.ConsumerMessage) error {
	cg.consumer.mu.RLock()
	handler, exists := cg.consumer.handlers[message.Type]
	cg.consumer.mu.RUnlock()

	if !exists {
		return fmt.Errorf("no handler registered for message type: %s", message.Type)
	}

	// Add correlation ID to context if present
	if message.Metadata.CorrelationID != "" {
		ctx = context.WithValue(ctx, "correlation_id", message.Metadata.CorrelationID)
	}

	// Process the message with retry logic
	return cg.processWithRetry(ctx, message, handler, saramaMessage)
}

// processWithRetry processes a message with retry logic
func (cg *ConsumerGroup) processWithRetry(ctx context.Context, message *domain.KafkaMessage, handler MessageHandler, saramaMessage *sarama.ConsumerMessage) error {
	maxRetries := cg.consumer.config.Kafka.RetryMaxAttempts

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := handler(ctx, message)
		if err == nil {
			return nil
		}

		cg.consumer.logger.WithFields(logrus.Fields{
			"message_id":   message.ID,
			"message_type": message.Type,
			"attempt":      attempt + 1,
			"max_retries":  maxRetries,
			"error":        err.Error(),
		}).Warn("Message processing failed, retrying...")

		// If this is the last attempt, send to dead letter queue
		if attempt == maxRetries {
			return cg.sendToDeadLetter(message, err, saramaMessage)
		}

		// Wait before retrying
		backoff := cg.consumer.config.Kafka.RetryBackoff * time.Duration(attempt+1)
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
func (cg *ConsumerGroup) sendToDeadLetter(message *domain.KafkaMessage, err error, saramaMessage *sarama.ConsumerMessage) error {
	// Parse message ID as UUID
	messageID, parseErr := uuid.Parse(message.ID)
	if parseErr != nil {
		cg.consumer.logger.WithError(parseErr).Error("Failed to parse message ID as UUID")
		messageID = uuid.New() // Fallback to new UUID
	}

	deadLetterMsg := &domain.DeadLetterMessage{
		ID:         messageID,
		Topic:      saramaMessage.Topic,
		Partition:  saramaMessage.Partition,
		Offset:     saramaMessage.Offset,
		Key:        saramaMessage.Key,
		Value:      saramaMessage.Value,
		Error:      err.Error(),
		RetryCount: message.Metadata.RetryCount,
		MaxRetries: cg.consumer.config.Kafka.RetryMaxAttempts,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Here you would typically save the dead letter message to the database
	// For now, we'll just log it
	cg.consumer.logger.WithFields(logrus.Fields{
		"message_id":  deadLetterMsg.ID,
		"topic":       deadLetterMsg.Topic,
		"error":       deadLetterMsg.Error,
		"retry_count": deadLetterMsg.RetryCount,
		"max_retries": deadLetterMsg.MaxRetries,
	}).Error("Message sent to dead letter queue")

	return err
}

// Close closes the consumer
func (c *Consumer) Close() error {
	c.cancel()
	return c.consumer.Close()
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
