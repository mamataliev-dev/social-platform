package logger

import (
	"log/slog"
	"os"
)

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
