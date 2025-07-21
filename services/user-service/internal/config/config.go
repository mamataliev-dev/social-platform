// Package config provides configuration loading and parsing for the user-service.
// It supports YAML-based config files and environment variable expansion, enabling
// flexible deployment and adherence to the Single Responsibility and Dependency
// Inversion principles.
package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the user-service, including server, gRPC,
// database, security, logging, and JWT settings. It is designed for extension
// without modification (Open/Closed Principle).
type Config struct {
	Env      string     `yaml:"env"`
	Server   Server     `yaml:"server"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Database Database   `yaml:"database"`
	Security Security   `yaml:"security"`
	Logging  Logging    `yaml:"logging"`
	JWT      JWT        `yaml:"jwt"`
}

// Server contains HTTP server configuration parameters.
type Server struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

// KeepaliveConfig defines gRPC keepalive settings.
type KeepaliveConfig struct {
	Time    time.Duration `yaml:"time"`
	Timeout time.Duration `yaml:"timeout"`
}

// TLSConfig holds TLS settings for secure gRPC communication.
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

// GRPCConfig contains gRPC server configuration parameters.
type GRPCConfig struct {
	Host                 string          `yaml:"host"`
	Port                 int             `yaml:"port"`
	MaxConcurrentStreams uint32          `yaml:"max_concurrent_streams"`
	Keepalive            KeepaliveConfig `yaml:"keepalive"`
	TLS                  TLSConfig       `yaml:"tls"`
}

// Database holds database connection configuration.
type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

// JWT contains JWT signing and expiry configuration.
type JWT struct {
	Secret       string `yaml:"secret"`
	ExpiresInMin int    `yaml:"expires_in_minutes"`
}

// Security holds security-related configuration, such as allowed CORS origins.
type Security struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

// Logging defines logging level and format configuration.
type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load reads and parses the configuration from the given YAML file path, expanding
// environment variables. It returns a Config struct or an error. This function
// enables Dependency Inversion by decoupling config source from consumers.
func Load(path string) (*Config, error) {
	var cfg Config

	if err := godotenv.Load(".env"); err == nil {
		slog.Info("Environment variables loaded from .env files")
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read config file", "err", err)
		return nil, err
	}

	expanded := os.ExpandEnv(string(raw))

	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		slog.Error("Failed to parse YAML config", "err", err)
		return nil, err
	}

	slog.Info("Config loaded successfully", "env", cfg.Env)
	return &cfg, nil
}
