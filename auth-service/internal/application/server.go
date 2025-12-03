package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router   *chi.Mux
	Handlers *Handlers
	Config   *config.Config
}

func NewServer(handlers *Handlers, cfg *config.Config) *Server {
	log.Println("Creating server")
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(corsMiddleware(cfg.CORS))

	return &Server{
		Router:   r,
		Handlers: handlers,
		Config:   cfg,
	}
}

func (s *Server) MountRoutes() {
	log.Println("Mounting routes")

	s.Router.Get("/health", s.Handlers.HandleHealth)

	s.Router.Post("/register", s.Handlers.HandleRegistration)
	s.Router.Post("/login", s.Handlers.HandleLogin)
	s.Router.Post("/refresh", s.Handlers.HandleRefreshToken)
	s.Router.Post("/logout", s.Handlers.HandleLogout)
	s.Router.Get("/me", s.Handlers.HandleMe)
	s.Router.Get("/yandex/login", s.Handlers.HandleYandexLogin)
	s.Router.Get("/google/login", s.Handlers.HandleGoogleLogin)
	s.Router.Get("/callback/yandex", s.Handlers.HandleYandexCallback)
	s.Router.Get("/callback/google", s.Handlers.HandleGoogleCallback)
	s.Router.Post("/password/reset", s.Handlers.HandlePasswordReset)
	s.Router.Post("/password/forgot", s.Handlers.HandleForgotPassword)
	s.Router.Get("/email/verify", s.Handlers.HandleEmailVerification)
	s.Router.Post("/email/resend", s.Handlers.HandleEmailResend)

	log.Println("Routes mounted")
}

func corsMiddleware(corsConfig config.CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			if !corsConfig.IsOriginAllowed(origin) {
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Set-Cookie")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (s *Server) Start(port string) error {
	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      s.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error, 1)

	go func() {
		log.Printf("Starting server on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case sig := <-shutdownCh:
		log.Printf("Received signal %v, shutting down gracefully...", sig)
	case err := <-errCh:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("Server stopped gracefully")
	return nil
}
