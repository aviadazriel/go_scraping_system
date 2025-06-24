package models

import (
	"fmt"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// Frequency represents a scraping frequency
type Frequency string

const (
	Frequency30Seconds Frequency = "30s"
	Frequency1Minute   Frequency = "1m"
	Frequency5Minutes  Frequency = "5m"
	Frequency15Minutes Frequency = "15m"
	Frequency30Minutes Frequency = "30m"
	Frequency1Hour     Frequency = "1h"
	Frequency6Hours    Frequency = "6h"
	Frequency12Hours   Frequency = "12h"
	Frequency1Day      Frequency = "1d"
	Frequency1Week     Frequency = "1w"
)

// ParseFrequency parses a frequency string into a time.Duration
func ParseFrequency(frequency string) (time.Duration, error) {
	switch Frequency(frequency) {
	case Frequency30Seconds:
		return 30 * time.Second, nil
	case Frequency1Minute:
		return 1 * time.Minute, nil
	case Frequency5Minutes:
		return 5 * time.Minute, nil
	case Frequency15Minutes:
		return 15 * time.Minute, nil
	case Frequency30Minutes:
		return 30 * time.Minute, nil
	case Frequency1Hour:
		return 1 * time.Hour, nil
	case Frequency6Hours:
		return 6 * time.Hour, nil
	case Frequency12Hours:
		return 12 * time.Hour, nil
	case Frequency1Day:
		return 24 * time.Hour, nil
	case Frequency1Week:
		return 7 * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("unsupported frequency: %s", frequency)
	}
}

// CalculateNextScrapeTime calculates the next scrape time based on frequency
func CalculateNextScrapeTime(frequency string, from time.Time) (time.Time, error) {
	duration, err := ParseFrequency(frequency)
	if err != nil {
		return time.Time{}, err
	}
	return from.Add(duration), nil
}

// IsValidFrequency checks if a frequency string is valid
func IsValidFrequency(frequency string) bool {
	_, err := ParseFrequency(frequency)
	return err == nil
}
