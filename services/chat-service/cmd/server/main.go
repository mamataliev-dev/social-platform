package main

import (
	"fmt"
	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/repository"
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

	// Setup Postgres connection + repository
	db, err := repository.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("Failed to connect to Postgres", "error", err)
		os.Exit(1)
	}

	fmt.Println(db)
}
