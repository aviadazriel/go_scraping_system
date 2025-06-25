// This file was moved from pkg/kafka/producer.go

package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

// Producer represents a Kafka producer using kafka-go
// (domain and config dependencies should be refactored to shared or injected)
type Producer struct {
	writers map[string]*kafka.Writer
	brokers []string
	logger  *logrus.Logger
	mu      sync.RWMutex
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, log *logrus.Logger) (*Producer, error) {
	return &Producer{
		writers: make(map[string]*kafka.Writer),
		brokers: brokers,
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
		Brokers:      p.brokers,
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
func (p *Producer) SendMessage(ctx context.Context, topic string, key string, value interface{}, headers map[string]string) error {
	writer := p.getWriter(topic)

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	var kafkaHeaders []kafka.Header
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafka.Header{Key: k, Value: []byte(v)})
	}

	kafkaMsg := kafka.Message{
		Key:     []byte(key),
		Value:   data,
		Headers: kafkaHeaders,
	}

	err = writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to send message to topic %s: %w", topic, err)
	}

	p.logger.WithFields(logrus.Fields{
		"topic": topic,
		"key":   key,
	}).Debug("Message sent successfully")

	return nil
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
