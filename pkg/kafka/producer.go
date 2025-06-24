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

// Producer represents a Kafka producer using kafka-go and implements domain.KafkaProducer
type Producer struct {
	writers map[string]*kafka.Writer
	config  *config.Config
	logger  *logrus.Logger
	mu      sync.RWMutex
}

// Ensure Producer implements domain.KafkaProducer
var _ domain.KafkaProducer = (*Producer)(nil)

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config, log *logrus.Logger) (*Producer, error) {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		config:  cfg,
		logger:  log,
	}, nil
}

// getWriter returns or creates a writer for the given topic
func (p *Producer) getWriter(topic string) *kafka.Writer {
	p.mu.RLock()
	writer, exists := p.writers[topic]
	p.mu.RUnlock()

	if exists {
		return writer
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if writer, exists = p.writers[topic]; exists {
		return writer
	}

	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      p.config.Kafka.Brokers,
		Topic:        topic,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		RequiredAcks: -1,    // Require all replicas to acknowledge
		Async:        false, // Use sync for reliability
		Logger: kafka.LoggerFunc(func(msg string, args ...interface{}) {
			p.logger.Debugf(msg, args...)
		}),
	})

	p.writers[topic] = writer
	return writer
}

// SendMessage sends a message to a Kafka topic
func (p *Producer) SendMessage(ctx context.Context, topic string, message *domain.KafkaMessage) error {
	writer := p.getWriter(topic)

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(message.ID),
		Value: data,
		Headers: []kafka.Header{
			{Key: "message_type", Value: []byte(string(message.Type))},
			{Key: "correlation_id", Value: []byte(message.Metadata.CorrelationID)},
			{Key: "timestamp", Value: []byte(message.Timestamp.Format(time.RFC3339))},
		},
	}

	err = writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	p.logger.WithFields(logrus.Fields{
		"topic":      topic,
		"message_id": message.ID,
		"type":       message.Type,
	}).Debug("Message sent successfully")

	return nil
}

// SendScrapingTask sends a scraping task message
func (p *Producer) SendScrapingTask(ctx context.Context, task *domain.ScrapingTask) error {
	correlationID := uuid.New().String()
	message := domain.NewScrapingTaskMessage(task, correlationID)

	return p.SendMessage(ctx, domain.TopicScrapingTasks, message)
}

// SendScrapedData sends a scraped data message
func (p *Producer) SendScrapedData(ctx context.Context, data *domain.ScrapedData, success bool, err string) error {
	correlationID := uuid.New().String()
	message := domain.NewScrapedDataMessage(data, success, err, correlationID)

	return p.SendMessage(ctx, domain.TopicScrapedData, message)
}

// SendParsedData sends a parsed data message
func (p *Producer) SendParsedData(ctx context.Context, data *domain.ParsedData) error {
	correlationID := uuid.New().String()
	message := domain.NewParsedDataMessage(data, correlationID)

	return p.SendMessage(ctx, domain.TopicParsedData, message)
}

// SendDeadLetter sends a dead letter message
func (p *Producer) SendDeadLetter(ctx context.Context, originalMessage *domain.KafkaMessage, err error, maxRetries int) error {
	message := domain.NewDeadLetterMessage(originalMessage, err, p.config.Kafka.RetryMaxAttempts)

	return p.SendMessage(ctx, domain.TopicDeadLetter, message)
}

// SendRetryMessage sends a retry message
func (p *Producer) SendRetryMessage(ctx context.Context, originalMessageID string, messageType domain.MessageType, data interface{}, retryCount, maxRetries int, retryDelay time.Duration) error {
	message := domain.NewRetryMessage(
		originalMessageID,
		messageType,
		data,
		retryCount,
		p.config.Kafka.RetryMaxAttempts,
		p.config.Kafka.RetryBackoff*time.Duration(retryCount+1),
	)

	return p.SendMessage(ctx, domain.TopicRetry, message)
}

// Close closes all writers
func (p *Producer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	var lastErr error
	for topic, writer := range p.writers {
		if err := writer.Close(); err != nil {
			p.logger.WithError(err).WithField("topic", topic).Error("Failed to close writer")
			lastErr = err
		}
	}

	return lastErr
}

// HealthCheck performs a producer health check
func (p *Producer) HealthCheck(ctx context.Context) error {
	// Send a test message to check if producer is working
	testMessage := &domain.KafkaMessage{
		ID:        uuid.New().String(),
		Type:      domain.MessageTypeScrapingTask,
		Timestamp: time.Now(),
		Data:      map[string]string{"test": "health_check"},
		Metadata: domain.Metadata{
			CorrelationID: uuid.New().String(),
			Source:        "health_check",
		},
	}

	return p.SendMessage(ctx, "health-check", testMessage)
}
