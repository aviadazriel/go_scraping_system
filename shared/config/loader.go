package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader handles loading configuration files with inheritance from shared configs
type Loader struct {
	viper *viper.Viper
}

// NewLoader creates a new configuration loader
func NewLoader() *Loader {
	return &Loader{
		viper: viper.New(),
	}
}

// LoadServiceConfig loads a service-specific configuration with inheritance from shared config
func (l *Loader) LoadServiceConfig(serviceName string) error {
	// Set up viper configuration
	l.viper.SetConfigName(serviceName)
	l.viper.SetConfigType("yaml")
	l.viper.AddConfigPath("../../configs") // Relative to shared/config/
	l.viper.AddConfigPath("../configs")    // Alternative path
	l.viper.AddConfigPath("./configs")     // Current directory

	// First, load shared configuration
	if err := l.loadSharedConfig(); err != nil {
		return fmt.Errorf("failed to load shared config: %w", err)
	}

	// Then, load service-specific configuration (will override shared values)
	if err := l.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read service config: %w", err)
	}

	return nil
}

// loadSharedConfig loads the shared configuration file
func (l *Loader) loadSharedConfig() error {
	sharedViper := viper.New()
	sharedViper.SetConfigName("shared")
	sharedViper.SetConfigType("yaml")
	sharedViper.AddConfigPath("../../configs")
	sharedViper.AddConfigPath("../configs")
	sharedViper.AddConfigPath("./configs")

	if err := sharedViper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read shared config: %w", err)
	}

	// Merge shared config into main viper instance
	for _, key := range sharedViper.AllKeys() {
		l.viper.Set(key, sharedViper.Get(key))
	}

	return nil
}

// GetString retrieves a string value from configuration
func (l *Loader) GetString(key string) string {
	return l.viper.GetString(key)
}

// GetInt retrieves an integer value from configuration
func (l *Loader) GetInt(key string) int {
	return l.viper.GetInt(key)
}

// GetBool retrieves a boolean value from configuration
func (l *Loader) GetBool(key string) bool {
	return l.viper.GetBool(key)
}

// GetStringSlice retrieves a string slice from configuration
func (l *Loader) GetStringSlice(key string) []string {
	return l.viper.GetStringSlice(key)
}

// GetDuration retrieves a duration value from configuration
func (l *Loader) GetDuration(key string) string {
	return l.viper.GetString(key)
}

// GetSub returns a sub-configuration for a given key
func (l *Loader) GetSub(key string) *viper.Viper {
	return l.viper.Sub(key)
}

// AllKeys returns all configuration keys
func (l *Loader) AllKeys() []string {
	return l.viper.AllKeys()
}

// IsSet checks if a key is set in configuration
func (l *Loader) IsSet(key string) bool {
	return l.viper.IsSet(key)
}

// LoadFromEnv loads configuration from environment variables
func (l *Loader) LoadFromEnv() {
	l.viper.AutomaticEnv()
	l.viper.SetEnvPrefix("SCRAPER")
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

// LoadFromFile loads configuration from a specific file
func (l *Loader) LoadFromFile(configPath string) error {
	l.viper.SetConfigFile(configPath)
	return l.viper.ReadInConfig()
}
