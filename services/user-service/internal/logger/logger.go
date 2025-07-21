// Package logger provides structured logging setup for the user-service.
// It configures slog for both development and production environments,
// following the Single Responsibility Principle.
package logger

import (
	"log/slog"
	"os"
)

// SetupLogger configures the global slog logger based on the environment.
// It enables Dependency Inversion by abstracting logger setup from consumers.
func SetupLogger(env string) {
	var handler slog.Handler

	switch env {
	case "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default: // development
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: false,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
