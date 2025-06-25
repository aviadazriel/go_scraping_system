package main

import (
	"fmt"
	"testing"

	"go_scraping_project/shared/config"
)

func TestConfigLoading(t *testing.T) {
	// Test that configuration can be loaded
	loader := config.NewLoader()
	err := loader.LoadServiceConfig("url-manager")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test that basic configuration values are loaded
	logLevel := loader.GetString("logging.level")
	if logLevel == "" {
		t.Log("Log level not set, using default")
	}

	dbHost := loader.GetString("database.host")
	if dbHost == "" {
		t.Fatal("Database host not configured")
	}

	kafkaBrokers := loader.GetStringSlice("kafka.brokers")
	if len(kafkaBrokers) == 0 {
		t.Fatal("Kafka brokers not configured")
	}

	t.Logf("Configuration loaded successfully:")
	t.Logf("  Database Host: %s", dbHost)
	t.Logf("  Kafka Brokers: %v", kafkaBrokers)
	t.Logf("  Log Level: %s", logLevel)
}

func TestDatabaseURLGeneration(t *testing.T) {
	loader := config.NewLoader()
	err := loader.LoadServiceConfig("url-manager")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test database URL generation
	dbHost := loader.GetString("database.host")
	dbPort := loader.GetInt("database.port")
	dbUser := loader.GetString("database.user")
	dbPassword := loader.GetString("database.password")
	dbName := loader.GetString("database.db_name")
	dbSSLMode := loader.GetString("database.ssl_mode")

	if dbName == "" {
		dbName = "scraping_db"
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	if databaseURL == "" {
		t.Fatal("Failed to generate database URL")
	}

	t.Logf("Generated database URL: %s", databaseURL)
}
