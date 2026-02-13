package middleware

import (
	"crypto/subtle"
	"net/http"
)

func APIKeyAuth(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			provided := r.Header.Get("X-API-Key")
			if provided == "" {
				http.Error(w, `{"error":"X-API-Key header is required"}`, http.StatusUnauthorized)
				return
			}

			if subtle.ConstantTimeCompare([]byte(provided), []byte(apiKey)) != 1 {
				http.Error(w, `{"error":"invalid API key"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
