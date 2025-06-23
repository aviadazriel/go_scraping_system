package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go_scraping_project/internal/config"
	"go_scraping_project/internal/domain"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	config   *config.Config
	logger   *logrus.Logger
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg *config.Config, log *logrus.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = cfg.Kafka.RetryMaxAttempts
	config.Producer.Retry.Backoff = cfg.Kafka.RetryBackoff
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Timeout = 10 * time.Second

	// Use the latest version
	config.Version = sarama.V2_8_0_0

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		producer: producer,
		config:   cfg,
		logger:   log,
	}, nil
}

// SendMessage sends a message to a Kafka topic
func (p *Producer) SendMessage(ctx context.Context, topic string, message *domain.KafkaMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.ID),
		Value: sarama.ByteEncoder(data),
		Headers: []sarama.RecordHeader{
			{Key: []byte("message_type"), Value: []byte(string(message.Type))},
			{Key: []byte("correlation_id"), Value: []byte(message.Metadata.CorrelationID)},
			{Key: []byte("timestamp"), Value: []byte(message.Timestamp.Format(time.RFC3339))},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	p.logger.WithFields(logrus.Fields{
		"topic":     topic,
		"partition": partition,
		"offset":    offset,
		"message_id": message.ID,
		"type":      message.Type,
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

// SendDeadLetterMessage sends a dead letter message
func (p *Producer) SendDeadLetterMessage(ctx context.Context, originalMessage *domain.KafkaMessage, err error) error {
	message := domain.NewDeadLetterMessage(originalMessage, err, p.config.Kafka.RetryMaxAttempts)

	return p.SendMessage(ctx, domain.TopicDeadLetter, message)
}

// SendRetryMessage sends a retry message
func (p *Producer) SendRetryMessage(ctx context.Context, originalMessageID string, messageType domain.MessageType, data interface{}, retryCount int) error {
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

// Close closes the producer
func (p *Producer) Close() error {
	return p.producer.Close()
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