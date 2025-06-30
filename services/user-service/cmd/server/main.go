package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/mamataliev-dev/social-platform/user-service/internal/config"
	"github.com/mamataliev-dev/social-platform/user-service/internal/logger"
)

func main() {
	// Load .ENV
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	// Setup router
	r := chi.NewRouter()

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	slog.Info("Starting server", "address", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
