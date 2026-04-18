package http

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type AnalysisHandler interface {
	HandleGetAnalysis(w http.ResponseWriter, r *http.Request)
	HandleGetAnalysesByTicker(w http.ResponseWriter, r *http.Request)
	HandleGetReportResults(w http.ResponseWriter, r *http.Request)
	HandleGetLatestReportResults(w http.ResponseWriter, r *http.Request)
	HandleGetBusinessResearch(w http.ResponseWriter, r *http.Request)
	HandleGetNews(w http.ResponseWriter, r *http.Request)
	HandleTriggerNews(w http.ResponseWriter, r *http.Request)
}

type HttpServer struct {
	srv             *http.Server
	analysisHandler AnalysisHandler
}

func NewHttpServer(analysisHandler AnalysisHandler) *HttpServer {
	return &HttpServer{
		analysisHandler: analysisHandler,
	}
}

func (h *HttpServer) RegisterRoutes(port int, apiKey string) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/analysis", h.analysisHandler.HandleGetAnalysis)
	r.Get("/analyses", h.analysisHandler.HandleGetAnalysesByTicker)
	r.Get("/report-results", h.analysisHandler.HandleGetReportResults)
	r.Get("/report-results/latest", h.analysisHandler.HandleGetLatestReportResults)
	r.Get("/business-research", h.analysisHandler.HandleGetBusinessResearch)
	r.Get("/news", h.analysisHandler.HandleGetNews)
	r.Post("/news/trigger", h.analysisHandler.HandleTriggerNews)

	r.Group(func(r chi.Router) {
		r.Use(apiKeyAuth(apiKey))
	})

	addr := fmt.Sprintf(":%d", port)
	slog.Info("AI Service starting", slog.String("addr", addr))

	h.srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func (h *HttpServer) RunServer(ctx context.Context) error {
	if h.srv == nil {
		return errors.New("server is no initialized")
	}

	errChan := make(chan error, 1)

	go func() {
		if err := h.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		slog.Error("Failed to start server", slog.Any("error", err))
		return err
	case <-ctx.Done():
		slog.Info("Shutting down gracefully...")
		h.srv.Shutdown(ctx)
	}

	close(errChan)

	return nil
}

func apiKeyAuth(apiKey string) func(http.Handler) http.Handler {
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
