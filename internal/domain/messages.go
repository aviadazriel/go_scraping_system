package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType represents the type of Kafka message
type MessageType string

const (
	MessageTypeScrapingTask MessageType = "scraping_task"
	MessageTypeScrapedData  MessageType = "scraped_data"
	MessageTypeParsedData   MessageType = "parsed_data"
	MessageTypeDeadLetter   MessageType = "dead_letter"
	MessageTypeRetry        MessageType = "retry"
)

// KafkaMessage represents a generic Kafka message
type KafkaMessage struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Metadata  Metadata    `json:"metadata"`
}

// Metadata represents message metadata
type Metadata struct {
	CorrelationID string            `json:"correlation_id"`
	Source        string            `json:"source"`
	RetryCount    int               `json:"retry_count"`
	Headers       map[string]string `json:"headers,omitempty"`
}

// ScrapingTaskMessage represents a scraping task message
type ScrapingTaskMessage struct {
	TaskID    uuid.UUID `json:"task_id"`
	URLID     uuid.UUID `json:"url_id"`
	URL       string    `json:"url"`
	UserAgent string    `json:"user_agent"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
	Attempt   int       `json:"attempt"`
}

// ScrapedDataMessage represents a scraped data message
type ScrapedDataMessage struct {
	TaskID        uuid.UUID `json:"task_id"`
	URLID         uuid.UUID `json:"url_id"`
	ScrapedDataID uuid.UUID `json:"scraped_data_id"`
	FilePath      string    `json:"file_path"`
	StatusCode    int       `json:"status_code"`
	ScrapedAt     time.Time `json:"scraped_at"`
	Success       bool      `json:"success"`
	Error         string    `json:"error,omitempty"`
}

// ParsedDataMessage represents a parsed data message
type ParsedDataMessage struct {
	URLID         uuid.UUID       `json:"url_id"`
	ScrapedDataID uuid.UUID       `json:"scraped_data_id"`
	ParsedDataID  uuid.UUID       `json:"parsed_data_id"`
	Schema        string          `json:"schema"`
	Data          json.RawMessage `json:"data"`
	CreatedAt     time.Time       `json:"created_at"`
}

// KafkaDeadLetterMessage represents a dead letter message
type KafkaDeadLetterMessage struct {
	OriginalMessage KafkaMessage `json:"original_message"`
	Error           string       `json:"error"`
	FailedAt        time.Time    `json:"failed_at"`
	RetryCount      int          `json:"retry_count"`
	MaxRetries      int          `json:"max_retries"`
	NextRetryAt     *time.Time   `json:"next_retry_at,omitempty"`
}

// RetryMessage represents a retry message
type RetryMessage struct {
	OriginalMessageID string      `json:"original_message_id"`
	MessageType       MessageType `json:"message_type"`
	Data              interface{} `json:"data"`
	RetryCount        int         `json:"retry_count"`
	MaxRetries        int         `json:"max_retries"`
	RetryAt           time.Time   `json:"retry_at"`
}

// Kafka Topics
const (
	TopicScrapingTasks = "scraping-tasks"
	TopicScrapedData   = "scraped-data"
	TopicParsedData    = "parsed-data"
	TopicDeadLetter    = "dead-letter"
	TopicRetry         = "retry"
)

// MessageHeaders represents Kafka message headers
type MessageHeaders map[string][]byte

// NewKafkaMessage creates a new Kafka message
func NewKafkaMessage(messageType MessageType, data interface{}, correlationID string) *KafkaMessage {
	return &KafkaMessage{
		ID:        uuid.New().String(),
		Type:      messageType,
		Timestamp: time.Now(),
		Data:      data,
		Metadata: Metadata{
			CorrelationID: correlationID,
			RetryCount:    0,
			Headers:       make(map[string]string),
		},
	}
}

// NewScrapingTaskMessage creates a new scraping task message
func NewScrapingTaskMessage(task *ScrapingTask, correlationID string) *KafkaMessage {
	return NewKafkaMessage(MessageTypeScrapingTask, ScrapingTaskMessage{
		TaskID:    task.ID,
		URLID:     task.URLID,
		URL:       task.URL,
		CreatedAt: task.CreatedAt,
		Attempt:   task.Attempt,
	}, correlationID)
}

// NewScrapedDataMessage creates a new scraped data message
func NewScrapedDataMessage(data *ScrapedData, success bool, err string, correlationID string) *KafkaMessage {
	return NewKafkaMessage(MessageTypeScrapedData, ScrapedDataMessage{
		TaskID:        data.TaskID,
		URLID:         data.URLID,
		ScrapedDataID: data.ID,
		FilePath:      data.FilePath,
		StatusCode:    data.StatusCode,
		ScrapedAt:     data.ScrapedAt,
		Success:       success,
		Error:         err,
	}, correlationID)
}

// NewParsedDataMessage creates a new parsed data message
func NewParsedDataMessage(data *ParsedData, correlationID string) *KafkaMessage {
	return NewKafkaMessage(MessageTypeParsedData, ParsedDataMessage{
		URLID:         data.URLID,
		ScrapedDataID: data.ScrapedDataID,
		ParsedDataID:  data.ID,
		Schema:        data.Schema,
		Data:          data.Data,
		CreatedAt:     data.CreatedAt,
	}, correlationID)
}

// NewDeadLetterMessage creates a new dead letter message
func NewDeadLetterMessage(originalMessage *KafkaMessage, err error, maxRetries int) *KafkaMessage {
	return NewKafkaMessage(MessageTypeDeadLetter, KafkaDeadLetterMessage{
		OriginalMessage: *originalMessage,
		Error:           err.Error(),
		FailedAt:        time.Now(),
		RetryCount:      0,
		MaxRetries:      maxRetries,
	}, originalMessage.Metadata.CorrelationID)
}

// NewRetryMessage creates a new retry message
func NewRetryMessage(originalMessageID string, messageType MessageType, data interface{}, retryCount, maxRetries int, retryDelay time.Duration) *KafkaMessage {
	return NewKafkaMessage(MessageTypeRetry, RetryMessage{
		OriginalMessageID: originalMessageID,
		MessageType:       messageType,
		Data:              data,
		RetryCount:        retryCount,
		MaxRetries:        maxRetries,
		RetryAt:           time.Now().Add(retryDelay),
	}, "")
}
