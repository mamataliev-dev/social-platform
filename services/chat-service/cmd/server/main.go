package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/config"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/logger"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.SetupLogger(cfg.Env)
	slog.Info("Logger initialized", "env", cfg.Env)
}
