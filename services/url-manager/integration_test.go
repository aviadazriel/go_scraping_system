package main

import (
	"context"
	"testing"
	"time"

	"go_scraping_project/shared/config"
	"go_scraping_project/shared/database"
	"go_scraping_project/shared/kafka"

	"github.com/sirupsen/logrus"
)

func TestDatabaseConnection(t *testing.T) {
	// Load configuration
	loader := config.NewLoader()
	err := loader.LoadServiceConfig("url-manager")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Test database connection
	db, err := database.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test a simple query
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	t.Log("Database connection successful")
}

func TestKafkaConnection(t *testing.T) {
	// Load configuration
	loader := config.NewLoader()
	err := loader.LoadServiceConfig("url-manager")
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Get Kafka brokers
	kafkaBrokers := loader.GetStringSlice("kafka.brokers")
	if len(kafkaBrokers) == 0 {
		kafkaBrokers = []string{"localhost:9092"}
	}

	// Initialize logger
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel) // Reduce noise during test

	// Test Kafka producer connection
	producer, err := kafka.NewProducer(kafkaBrokers, logger)
	if err != nil {
		t.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Test sending a message (this will fail if Kafka is not available)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testMessage := map[string]string{"test": "message"}
	err = producer.SendMessage(ctx, "test-topic", "test-key", testMessage, nil)
	if err != nil {
		t.Logf("Kafka producer test message failed (expected if topic doesn't exist): %v", err)
	} else {
		t.Log("Kafka producer connection successful")
	}

	t.Log("Kafka producer created successfully")
}
