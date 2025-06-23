package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Scraping ScrapingConfig `mapstructure:"scraping"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Metrics  MetricsConfig  `mapstructure:"metrics"`
	Tracing  TracingConfig  `mapstructure:"tracing"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	MigrationPath   string        `mapstructure:"migration_path"`
}

// KafkaConfig represents Kafka configuration
type KafkaConfig struct {
	Brokers           []string      `mapstructure:"brokers"`
	GroupID           string        `mapstructure:"group_id"`
	AutoOffsetReset   string        `mapstructure:"auto_offset_reset"`
	EnableAutoCommit  bool          `mapstructure:"enable_auto_commit"`
	AutoCommitInterval time.Duration `mapstructure:"auto_commit_interval"`
	SessionTimeout    time.Duration `mapstructure:"session_timeout"`
	HeartbeatInterval time.Duration `mapstructure:"heartbeat_interval"`
	MaxPollRecords    int           `mapstructure:"max_poll_records"`
	MaxPollInterval   time.Duration `mapstructure:"max_poll_interval"`
	RetryBackoff      time.Duration `mapstructure:"retry_backoff"`
	RetryMaxAttempts  int           `mapstructure:"retry_max_attempts"`
}

// ScrapingConfig represents scraping configuration
type ScrapingConfig struct {
	DefaultTimeout     time.Duration `mapstructure:"default_timeout"`
	DefaultUserAgent   string        `mapstructure:"default_user_agent"`
	DefaultRateLimit   int           `mapstructure:"default_rate_limit"`
	MaxRetries         int           `mapstructure:"max_retries"`
	RetryDelay         time.Duration `mapstructure:"retry_delay"`
	HTMLStoragePath    string        `mapstructure:"html_storage_path"`
	MaxConcurrentTasks int           `mapstructure:"max_concurrent_tasks"`
	RespectRobotsTxt   bool          `mapstructure:"respect_robots_txt"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	IncludeCaller bool `mapstructure:"include_caller"`
}

// MetricsConfig represents metrics configuration
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    int    `mapstructure:"port"`
	Path    string `mapstructure:"path"`
}

// TracingConfig represents tracing configuration
type TracingConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ServiceName string `mapstructure:"service_name"`
	JaegerURL   string `mapstructure:"jaeger_url"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Set defaults
	setDefaults()

	// Read environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SCRAPING")

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables
	overrideWithEnvVars()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "60s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	viper.SetDefault("database.migration_path", "./migrations")

	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.group_id", "scraper-group")
	viper.SetDefault("kafka.auto_offset_reset", "earliest")
	viper.SetDefault("kafka.enable_auto_commit", true)
	viper.SetDefault("kafka.auto_commit_interval", "1s")
	viper.SetDefault("kafka.session_timeout", "30s")
	viper.SetDefault("kafka.heartbeat_interval", "3s")
	viper.SetDefault("kafka.max_poll_records", 500)
	viper.SetDefault("kafka.max_poll_interval", "5m")
	viper.SetDefault("kafka.retry_backoff", "100ms")
	viper.SetDefault("kafka.retry_max_attempts", 3)

	// Scraping defaults
	viper.SetDefault("scraping.default_timeout", "30s")
	viper.SetDefault("scraping.default_user_agent", "GoScraper/1.0")
	viper.SetDefault("scraping.default_rate_limit", 1)
	viper.SetDefault("scraping.max_retries", 3)
	viper.SetDefault("scraping.retry_delay", "5s")
	viper.SetDefault("scraping.html_storage_path", "./data/html")
	viper.SetDefault("scraping.max_concurrent_tasks", 10)
	viper.SetDefault("scraping.respect_robots_txt", true)

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.include_caller", false)

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.port", 9090)
	viper.SetDefault("metrics.path", "/metrics")

	// Tracing defaults
	viper.SetDefault("tracing.enabled", false)
	viper.SetDefault("tracing.service_name", "scraper-service")
	viper.SetDefault("tracing.jaeger_url", "http://localhost:14268/api/traces")
}

// overrideWithEnvVars overrides configuration with environment variables
func overrideWithEnvVars() {
	// Server
	if port := os.Getenv("SCRAPING_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			viper.Set("server.port", p)
		}
	}

	// Database
	if host := os.Getenv("SCRAPING_DATABASE_HOST"); host != "" {
		viper.Set("database.host", host)
	}
	if port := os.Getenv("SCRAPING_DATABASE_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			viper.Set("database.port", p)
		}
	}
	if user := os.Getenv("SCRAPING_DATABASE_USER"); user != "" {
		viper.Set("database.user", user)
	}
	if password := os.Getenv("SCRAPING_DATABASE_PASSWORD"); password != "" {
		viper.Set("database.password", password)
	}
	if database := os.Getenv("SCRAPING_DATABASE_NAME"); database != "" {
		viper.Set("database.database", database)
	}

	// Kafka
	if brokers := os.Getenv("SCRAPING_KAFKA_BROKERS"); brokers != "" {
		viper.Set("kafka.brokers", []string{brokers})
	}
	if groupID := os.Getenv("SCRAPING_KAFKA_GROUP_ID"); groupID != "" {
		viper.Set("kafka.group_id", groupID)
	}

	// Logging
	if level := os.Getenv("SCRAPING_LOG_LEVEL"); level != "" {
		viper.Set("logging.level", level)
	}
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 {
		return fmt.Errorf("server port must be positive")
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Port <= 0 {
		return fmt.Errorf("database port must be positive")
	}

	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if len(config.Kafka.Brokers) == 0 {
		return fmt.Errorf("at least one Kafka broker is required")
	}

	if config.Kafka.GroupID == "" {
		return fmt.Errorf("Kafka group ID is required")
	}

	if config.Scraping.MaxConcurrentTasks <= 0 {
		return fmt.Errorf("max concurrent tasks must be positive")
	}

	return nil
}

// GetDatabaseURL returns the database connection URL
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.SSLMode,
	)
} 