package utils

import (
	"fmt"
	"time"
)

// Time utilities for consistent timezone handling across services

// Now returns the current time in UTC
func Now() time.Time {
	return time.Now().UTC()
}

// FormatTime formats a time.Time to RFC3339 format in UTC
func FormatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ParseTime parses a time string in RFC3339 format and returns UTC time
func ParseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}
	return t.UTC(), nil
}

// ParseDuration parses a duration string (e.g., "1h", "30m", "1d")
func ParseDuration(durationStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format: %w", err)
	}
	return duration, nil
}

// CalculateNextScrapeTime calculates the next scrape time based on frequency
func CalculateNextScrapeTime(frequency string, from time.Time) (time.Time, error) {
	duration, err := ParseDuration(frequency)
	if err != nil {
		return time.Time{}, err
	}
	return from.Add(duration), nil
}

// IsTimeInRange checks if a time is within the specified range
func IsTimeInRange(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

// TimeRange represents a time range with start and end times
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewTimeRange creates a new time range
func NewTimeRange(start, end time.Time) TimeRange {
	return TimeRange{
		Start: start.UTC(),
		End:   end.UTC(),
	}
}

// Contains checks if the time range contains the given time
func (tr TimeRange) Contains(t time.Time) bool {
	return IsTimeInRange(t, tr.Start, tr.End)
}

// Duration returns the duration of the time range
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}
