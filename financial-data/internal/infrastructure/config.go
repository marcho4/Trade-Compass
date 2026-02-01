package infrastructure

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	URL               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type SecurityConfig struct {
	AdminAPIKey    string
	AllowedOrigins string
}

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	return &Config{
		Server: ServerConfig{
			Port:         getEnvOrDefault("SERVER_PORT", "8082"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			URL:               dbURL,
			MaxConns:          getEnvAsInt32("DB_MAX_CONNS", 25),
			MinConns:          getEnvAsInt32("DB_MIN_CONNS", 5),
			MaxConnLifetime:   getDurationEnv("DB_MAX_CONN_LIFETIME", 30*time.Minute),
			MaxConnIdleTime:   getDurationEnv("DB_MAX_CONN_IDLE_TIME", 5*time.Minute),
			HealthCheckPeriod: time.Minute,
		},
		Security: SecurityConfig{
			AdminAPIKey:    os.Getenv("ADMIN_API_KEY"),
			AllowedOrigins: getEnvOrDefault("ALLOWED_ORIGINS", "*"),
		},
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt32(key string, defaultValue int) int32 {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return int32(defaultValue)
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return int32(defaultValue)
	}
	return int32(value)
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
