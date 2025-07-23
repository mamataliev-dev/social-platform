package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the user-service, including server, gRPC,
// database, security, logging, and JWT settings.
type Config struct {
	Env      string     `yaml:"env"`      // Application environment (e.g., "development", "production")
	Server   Server     `yaml:"server"`   // HTTP server settings
	GRPC     GRPCConfig `yaml:"grpc"`     // gRPC server settings
	Database Database   `yaml:"database"` // Database connection settings
	Security Security   `yaml:"security"` // Security-related settings (e.g., CORS)
	Logging  Logging    `yaml:"logging"`  // Logging level and format
	JWT      JWT        `yaml:"jwt"`      // JWT signing
}

// Server contains HTTP server configuration parameters.
type Server struct {
	Host  string `yaml:"host"`  // Server host address
	Port  string `yaml:"port"`  // Server port
	Debug bool   `yaml:"debug"` // Debug mode flag
}

// KeepaliveConfig defines gRPC keepalive settings for connection health checks.
type KeepaliveConfig struct {
	Time    time.Duration `yaml:"time"`    // Interval between pings
	Timeout time.Duration `yaml:"timeout"` // Timeout for ping ack
}

// TLSConfig holds TLS settings for secure gRPC communication.
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`   // Enable TLS encryption
	CertFile string `yaml:"cert_file"` // Path to TLS certificate file
	KeyFile  string `yaml:"key_file"`  // Path to TLS key file
}

// GRPCConfig contains gRPC server configuration parameters.
type GRPCConfig struct {
	Host                 string          `yaml:"host"`                   // gRPC server host
	Port                 int             `yaml:"port"`                   // gRPC server port
	MaxConcurrentStreams uint32          `yaml:"max_concurrent_streams"` // Max parallel RPC streams
	Keepalive            KeepaliveConfig `yaml:"keepalive"`              // Keepalive settings
	TLS                  TLSConfig       `yaml:"tls"`                    // TLS configuration
}

// Database holds database connection configuration.
type Database struct {
	Driver   string `yaml:"driver"`   // Database driver (e.g., "postgres")
	Host     string `yaml:"host"`     // Database host address
	Port     int    `yaml:"port"`     // Database port
	User     string `yaml:"user"`     // Database user
	Password string `yaml:"password"` // Database password
	Name     string `yaml:"name"`     // Database name
	SSLMode  string `yaml:"sslmode"`  // SSL mode (e.g., "disable", "require")
}

// JWT contains JWT signing secret and token expiry settings.
type JWT struct {
	Secret string `yaml:"secret"` // Secret key for signing JWT tokens
}

// Security holds security-related configuration, such as allowed CORS origins.
type Security struct {
	AllowedOrigins []string `yaml:"allowed_origins"` // List of allowed CORS origins
}

// Logging defines logging level and format configuration.
type Logging struct {
	Level  string `yaml:"level"`  // Log level (e.g., "info", "debug")
	Format string `yaml:"format"` // Log output format (e.g., "json", "text")
}

// Load reads and parses the YAML configuration from the specified file path.
// It loads environment variables from a .env file, expands them in the YAML,
// and unmarshals into a Config struct. Returns an error on failure.
func Load(path string) (*Config, error) {
	if err := godotenv.Load(".env"); err == nil {
		slog.Info("Loaded environment variables from .env file")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read config file", "error", err)
		return nil, err
	}

	expanded := os.ExpandEnv(string(raw))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		slog.Error("Failed to parse YAML config", "error", err)
		return nil, err
	}

	slog.Info("Configuration loaded successfully", "environment", cfg.Env)
	return &cfg, nil
}
