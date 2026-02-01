package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	contextKeyAPIKey = contextKey("api_key")
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			RespondWithError(w, r, http.StatusUnauthorized, "X-API-Key header is required", nil)
			return
		}

		expectedKey := os.Getenv("ADMIN_API_KEY")
		if expectedKey == "" {
			RespondWithError(w, r, http.StatusInternalServerError, "Server configuration error", fmt.Errorf("ADMIN_API_KEY not configured"))
			return
		}

		if apiKey != expectedKey {
			RespondWithError(w, r, http.StatusUnauthorized, "Invalid API key", nil)
			return
		}

		ctx := context.WithValue(r.Context(), contextKeyAPIKey, apiKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
		if allowedOrigins == "" {
			allowedOrigins = "*"
		}

		origin := r.Header.Get("Origin")
		if origin != "" && (allowedOrigins == "*" || contains(strings.Split(allowedOrigins, ","), origin)) {
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

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}

func validateApiKey(r *http.Request) (bool, error) {
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == "" {
		return false, fmt.Errorf("X-API-Key header is required")
	}

	expectedKey := os.Getenv("ADMIN_API_KEY")
	if expectedKey == "" {
		return false, fmt.Errorf("ADMIN_API_KEY not configured")
	}

	if apiKey != expectedKey {
		return false, fmt.Errorf("Invalid API key")
	}

	return true, nil
}
