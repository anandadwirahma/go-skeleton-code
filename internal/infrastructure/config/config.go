// Package config handles application configuration loading from environment variables.
// It uses godotenv to load .env files and populates a strongly-typed Config struct.
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	// Server
	AppEnv     string
	ServerPort string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBMaxOpen  int
	DBMaxIdle  int

	// Logger
	LogLevel string

	// External HTTP Client
	ExternalAPIBaseURL string
	ExternalAPITimeout int // in seconds
}

// Load reads the .env file (if present) and populates a Config struct.
// It does NOT fail if .env is missing — useful for containerised deployments
// where env vars are injected directly.
func Load(envFile string) (*Config, error) {
	// Attempt to load .env; ignore error if file doesn't exist.
	_ = godotenv.Load(envFile)

	dbMaxOpen, err := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	if err != nil {
		return nil, fmt.Errorf("config: invalid DB_MAX_OPEN_CONNS: %w", err)
	}
	dbMaxIdle, err := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	if err != nil {
		return nil, fmt.Errorf("config: invalid DB_MAX_IDLE_CONNS: %w", err)
	}
	extTimeout, err := strconv.Atoi(getEnv("EXTERNAL_API_TIMEOUT_SEC", "10"))
	if err != nil {
		return nil, fmt.Errorf("config: invalid EXTERNAL_API_TIMEOUT_SEC: %w", err)
	}

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "skeleton_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBMaxOpen:  dbMaxOpen,
		DBMaxIdle:  dbMaxIdle,

		LogLevel: getEnv("LOG_LEVEL", "info"),

		ExternalAPIBaseURL: getEnv("EXTERNAL_API_BASE_URL", "https://jsonplaceholder.typicode.com"),
		ExternalAPITimeout: extTimeout,
	}

	return cfg, nil
}

// DSN builds the PostgreSQL connection string from config fields.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

// IsProd returns true when running in production environment.
func (c *Config) IsProd() bool {
	return c.AppEnv == "production"
}

// getEnv returns the value of an environment variable or a default fallback.
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
