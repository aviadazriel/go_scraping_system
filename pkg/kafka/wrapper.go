package kafka

import (
	"context"

	"go_scraping_project/internal/domain"
)

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
