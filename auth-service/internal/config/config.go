package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	OAuth    OAuthConfig
	CORS     CORSConfig
	Frontend FrontendConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func (c *JWTConfig) AccessTokenTTLSeconds() int64 {
	return int64(c.AccessTokenTTL.Seconds())
}

type OAuthConfig struct {
	Yandex YandexOAuthConfig
	// Google GoogleOAuthConfig // TODO: добавить когда будет нужно
}

type YandexOAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type CORSConfig struct {
	AllowedOrigins []string
	AllowAll       bool
}

type FrontendConfig struct {
	URL          string
	CookieDomain string
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.Server.Port = getEnvOrDefault("PORT", "8080")

	cfg.Database.URL = os.Getenv("DB_URL")
	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}
	if len(cfg.JWT.Secret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 bytes")
	}
	cfg.JWT.AccessTokenTTL = 15 * time.Minute
	cfg.JWT.RefreshTokenTTL = 15 * 24 * time.Hour

	cfg.OAuth.Yandex.ClientID = os.Getenv("YANDEX_CLIENT_ID")
	cfg.OAuth.Yandex.ClientSecret = os.Getenv("YANDEX_CLIENT_SECRET")
	if cfg.OAuth.Yandex.ClientID == "" || cfg.OAuth.Yandex.ClientSecret == "" {
		return nil, fmt.Errorf("YANDEX_CLIENT_ID and YANDEX_CLIENT_SECRET environment variables are required")
	}

	if err := cfg.loadCORSConfig(); err != nil {
		return nil, fmt.Errorf("load CORS config: %w", err)
	}

	cfg.Frontend.URL = os.Getenv("FRONTEND_URL")
	cfg.Frontend.CookieDomain = os.Getenv("COOKIE_DOMAIN")

	return cfg, nil
}

func (c *Config) loadCORSConfig() error {
	allowAll := os.Getenv("CORS_ALLOW_ALL")
	if allowAll == "true" {
		c.CORS.AllowAll = true
		return nil
	}

	originsStr := os.Getenv("CORS_ALLOWED_ORIGINS")
	if originsStr == "" {
		return fmt.Errorf("CORS_ALLOWED_ORIGINS environment variable is required (or set CORS_ALLOW_ALL=true for development)")
	}

	origins := strings.Split(originsStr, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	c.CORS.AllowedOrigins = origins
	return nil
}

func (c *CORSConfig) IsOriginAllowed(origin string) bool {
	if c.AllowAll {
		return true
	}

	for _, allowed := range c.AllowedOrigins {
		if allowed == origin {
			return true
		}
	}

	return false
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
