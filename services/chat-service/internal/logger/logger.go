package logger

import (
	"log/slog"
	"os"
)

// SetupLogger configures a base *slog.Logger with your env and app name,
// registers it as the global default, and returns it.
//
//   - env:     "production" for JSON/Info, anything else for text/Debug.
func SetupLogger(env string) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{}

	switch env {
	case "production":
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		opts.Level = slog.LevelDebug
		opts.AddSource = false
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	base := slog.New(handler).With("app", "chat-service")

	slog.SetDefault(base)
	return base
}
