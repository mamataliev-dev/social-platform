package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string     `yaml:"env"`
	Server   Server     `yaml:"server"`
	GRPC     GRPCConfig `yaml:"grpc"`
	Database Database   `yaml:"database"`
	Security Security   `yaml:"security"`
	Logging  Logging    `yaml:"logging"`
	JWT      JWT        `yaml:"jwt"`
}

type Server struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Debug bool   `yaml:"debug"`
}

type KeepaliveConfig struct {
	Time    time.Duration `yaml:"time"`
	Timeout time.Duration `yaml:"timeout"`
}

type TLSConfig struct {
	Enabled  bool   `yaml:"enabled"`
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type GRPCConfig struct {
	Host                 string          `yaml:"host"`
	Port                 int             `yaml:"port"`
	MaxConcurrentStreams uint32          `yaml:"max_concurrent_streams"`
	Keepalive            KeepaliveConfig `yaml:"keepalive"`
	TLS                  TLSConfig       `yaml:"tls"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	SSLMode  string `yaml:"sslmode"`
}

type JWT struct {
	Secret       string `yaml:"secret"`
	ExpiresInMin int    `yaml:"expires_in_minutes"`
}

type Security struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
}

type Logging struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

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
