package kafka

import (
	"context"
	"time"

	"go_scraping_project/internal/domain"
)

// KafkaProducerWrapper wraps the kafka.Producer to implement domain.KafkaProducer
// This allows services to use the domain interface while using the concrete kafka.Producer
type KafkaProducerWrapper struct {
	producer *Producer
}

// NewKafkaProducerWrapper creates a new wrapper around the kafka.Producer
func NewKafkaProducerWrapper(producer *Producer) domain.KafkaProducer {
	return &KafkaProducerWrapper{
		producer: producer,
	}
}

// SendMessage sends a message to a Kafka topic
func (w *KafkaProducerWrapper) SendMessage(ctx context.Context, topic string, message *domain.KafkaMessage) error {
	return w.producer.SendMessage(ctx, topic, message)
}

// SendScrapingTask sends a scraping task message
func (w *KafkaProducerWrapper) SendScrapingTask(ctx context.Context, task *domain.ScrapingTask) error {
	return w.producer.SendScrapingTask(ctx, task)
}

// SendScrapedData sends a scraped data message
func (w *KafkaProducerWrapper) SendScrapedData(ctx context.Context, data *domain.ScrapedData, success bool, err string) error {
	return w.producer.SendScrapedData(ctx, data, success, err)
}

// SendParsedData sends a parsed data message
func (w *KafkaProducerWrapper) SendParsedData(ctx context.Context, data *domain.ParsedData) error {
	return w.producer.SendParsedData(ctx, data)
}

// SendDeadLetter sends a dead letter message
func (w *KafkaProducerWrapper) SendDeadLetter(ctx context.Context, originalMessage *domain.KafkaMessage, err error, maxRetries int) error {
	return w.producer.SendDeadLetterMessage(ctx, originalMessage, err)
}

// SendRetryMessage sends a retry message
func (w *KafkaProducerWrapper) SendRetryMessage(ctx context.Context, originalMessageID string, messageType domain.MessageType, data interface{}, retryCount, maxRetries int, retryDelay time.Duration) error {
	return w.producer.SendRetryMessage(ctx, originalMessageID, messageType, data, retryCount)
}

// Close closes the producer
func (w *KafkaProducerWrapper) Close() error {
	return w.producer.Close()
}

// KafkaConsumerWrapper wraps the kafka.Consumer to provide a cleaner interface
// This allows services to use a simplified consumer interface
type KafkaConsumerWrapper struct {
	consumer *Consumer
}

// NewKafkaConsumerWrapper creates a new wrapper around the kafka.Consumer
func NewKafkaConsumerWrapper(consumer *Consumer) *KafkaConsumerWrapper {
	return &KafkaConsumerWrapper{
		consumer: consumer,
	}
}

// RegisterHandler registers a message handler for a specific message type
func (w *KafkaConsumerWrapper) RegisterHandler(messageType domain.MessageType, handler func(ctx context.Context, message *domain.KafkaMessage) error) {
	w.consumer.RegisterHandler(messageType, handler)
}

// Consume starts consuming messages from the specified topics
func (w *KafkaConsumerWrapper) Consume(topics []string) error {
	return w.consumer.Consume(topics)
}

// Close closes the consumer
func (w *KafkaConsumerWrapper) Close() error {
	return w.consumer.Close()
}

// HealthCheck performs a consumer health check
func (w *KafkaConsumerWrapper) HealthCheck(ctx context.Context) error {
	return w.consumer.HealthCheck(ctx)
}
