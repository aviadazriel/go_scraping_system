package observability

import (
	"context"
	"os"

	"go_scraping_project/internal/config"

	"github.com/sirupsen/logrus"
)

// Logger represents a structured logger
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new structured logger
func NewLogger(cfg *config.Config) *Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Logging.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Set log format
	switch cfg.Logging.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}

	// Set output
	switch cfg.Logging.Output {
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "stderr":
		logger.SetOutput(os.Stderr)
	default:
		logger.SetOutput(os.Stdout)
	}

	// Include caller information if configured
	if cfg.Logging.IncludeCaller {
		logger.SetReportCaller(true)
	}

	return &Logger{
		logger: logger,
	}
}

// WithContext creates a logger with context information
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.logger.WithContext(ctx)

	// Add correlation ID if present in context
	if correlationID, ok := ctx.Value("correlation_id").(string); ok && correlationID != "" {
		entry = entry.WithField("correlation_id", correlationID)
	}

	// Add request ID if present in context
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		entry = entry.WithField("request_id", requestID)
	}

	return entry
}

// WithFields creates a logger with additional fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

// WithField creates a logger with a single field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

// WithError creates a logger with an error
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.logger.WithError(err)
}

// Debug logs a debug message
func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Panic logs a panic message
func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Panicf logs a formatted panic message
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// GetLogger returns the underlying logrus logger
func (l *Logger) GetLogger() *logrus.Logger {
	return l.logger
} 