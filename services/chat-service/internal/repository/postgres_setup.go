package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/mamataliev-dev/social-platform/services/chat-service/internal/config"
)

// NewPostgresConnection initializes and verifies a PostgreSQL connection.
// It constructs the DSN from cfg, opens the connection, pings to ensure availability,
// and logs success. Returns the *sql.DB or an error if any step fails.
func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return db, nil
}
