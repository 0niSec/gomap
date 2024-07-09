package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

// init initializes the logger
func init() {
	once.Do(func() {
		// Create a JSON handler for structured logging
		handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug, // Set default level to Info
		})

		logger = slog.New(handler)
	})
}

// GetLogger returns the singleton logger instance
func GetLogger() *slog.Logger {
	return logger
}

// SetLevel sets the logging level
func SetLevel(level slog.Level) {
	handler := logger.Handler()
	if _, ok := handler.(*slog.JSONHandler); ok {
		newHandler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		logger = slog.New(newHandler)
	}
}

// Debug logs a debug message
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
