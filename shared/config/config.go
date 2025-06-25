package config

import (
	"time"
)

// Config represents the base configuration structure
type Config struct {
	Environment string         `json:"environment"`
	LogLevel    string         `json:"log_level"`
	Database    DatabaseConfig `json:"database"`
	Kafka       KafkaConfig    `json:"kafka"`
	HTTP        HTTPConfig     `json:"http"`
	Scraping    ScrapingConfig `json:"scraping"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
	MaxConns int    `json:"max_conns"`
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Brokers        []string       `json:"brokers"`
	Topics         TopicsConfig   `json:"topics"`
	ConsumerGroup  string         `json:"consumer_group"`
	ProducerConfig ProducerConfig `json:"producer_config"`
	ConsumerConfig ConsumerConfig `json:"consumer_config"`
}

// TopicsConfig represents Kafka topics configuration
type TopicsConfig struct {
	ScrapingTasks string `json:"scraping_tasks"`
	ScrapedData   string `json:"scraped_data"`
	ParsedData    string `json:"parsed_data"`
	DeadLetter    string `json:"dead_letter"`
}

// ProducerConfig represents Kafka producer configuration
type ProducerConfig struct {
	RequiredAcks int           `json:"required_acks"`
	Timeout      time.Duration `json:"timeout"`
	Retries      int           `json:"retries"`
}

// ConsumerConfig represents Kafka consumer configuration
type ConsumerConfig struct {
	AutoOffsetReset   string        `json:"auto_offset_reset"`
	SessionTimeout    time.Duration `json:"session_timeout"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval"`
}

// HTTPConfig represents HTTP server configuration
type HTTPConfig struct {
	Port         int           `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}

// ScrapingConfig represents scraping configuration
type ScrapingConfig struct {
	DefaultTimeout    time.Duration `json:"default_timeout"`
	DefaultUserAgent  string        `json:"default_user_agent"`
	DefaultMaxRetries int           `json:"default_max_retries"`
	DefaultRateLimit  int           `json:"default_rate_limit"`
	Concurrency       int           `json:"concurrency"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Environment: "development",
		LogLevel:    "info",
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "scraper",
			Password: "scraper",
			DBName:   "scraping_db",
			SSLMode:  "disable",
			MaxConns: 10,
		},
		Kafka: KafkaConfig{
			Brokers: []string{"localhost:9092"},
			Topics: TopicsConfig{
				ScrapingTasks: "scraping-tasks",
				ScrapedData:   "scraped-data",
				ParsedData:    "parsed-data",
				DeadLetter:    "dead-letter",
			},
			ConsumerGroup: "scraping-group",
			ProducerConfig: ProducerConfig{
				RequiredAcks: 1,
				Timeout:      30 * time.Second,
				Retries:      3,
			},
			ConsumerConfig: ConsumerConfig{
				AutoOffsetReset:   "earliest",
				SessionTimeout:    30 * time.Second,
				HeartbeatInterval: 3 * time.Second,
			},
		},
		HTTP: HTTPConfig{
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		Scraping: ScrapingConfig{
			DefaultTimeout:    30 * time.Second,
			DefaultUserAgent:  "GoScrapingBot/1.0",
			DefaultMaxRetries: 3,
			DefaultRateLimit:  1,
			Concurrency:       10,
		},
	}
}
