package middleware

import (
	"context"
	"errors"
	"financial_data/internal/application/response"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	contextKeyAPIKey = contextKey("api_key")
)

type MiddlewareConfig struct {
	AdminAPIKey    string
	allowedOrigins map[string]struct{}
	allowAll       bool
}

func NewMiddlewareConfig() (*MiddlewareConfig, error) {
	apiKey := os.Getenv("ADMIN_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ADMIN_API_KEY is not set")
	}

	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")

	config := &MiddlewareConfig{
		AdminAPIKey:    apiKey,
		allowedOrigins: make(map[string]struct{}),
	}

	if allowedOriginsStr == "" || allowedOriginsStr == "*" {
		config.allowAll = true
	} else {
		for _, origin := range strings.Split(allowedOriginsStr, ",") {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				config.allowedOrigins[trimmed] = struct{}{}
			}
		}
	}

	return config, nil
}

func (m *MiddlewareConfig) isOriginAllowed(origin string) bool {
	if m.allowAll {
		return true
	}
	_, exists := m.allowedOrigins[origin]
	return exists
}

func (m *MiddlewareConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			response.RespondWithError(w, r, http.StatusUnauthorized, "X-API-Key header is required", nil)
			return
		}

		if apiKey != m.AdminAPIKey {
			response.RespondWithError(w, r, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAPIKey, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareConfig) CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && m.isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareConfig) ValidateApiKey(r *http.Request) (bool, error) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		return false, fmt.Errorf("X-API-Key header is required")
	}

	if apiKey != m.AdminAPIKey {
		return false, fmt.Errorf("invalid API key")
	}

	return true, nil
}
