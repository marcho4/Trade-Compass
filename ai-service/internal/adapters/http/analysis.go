package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"ai-service/internal/domain"
	"ai-service/internal/domain/entity"
)

type analysisReader interface {
	GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error)
	GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error)
}

type reportResultsReader interface {
	GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error)
	GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error)
}

type businessResearchReader interface {
	GetBusinessResearch(ctx context.Context, ticker string) (*entity.BusinessResearchResult, error)
}

type newsReader interface {
	GetFreshNews(ctx context.Context, ticker string, ttl time.Duration) (*entity.NewsResponse, error)
}

type newsPublisher interface {
	PublishMessage(ctx context.Context, value []byte) error
}

type analysisHandler struct {
	analysis         analysisReader
	reportResults    reportResultsReader
	businessResearch businessResearchReader
	news             newsReader
	newsPublisher    newsPublisher
}

func NewAnalysisHandler(
	analysis analysisReader,
	reportResults reportResultsReader,
	businessResearch businessResearchReader,
	news newsReader,
	newsPublisher newsPublisher,
) *analysisHandler {
	return &analysisHandler{
		analysis:         analysis,
		reportResults:    reportResults,
		businessResearch: businessResearch,
		news:             news,
		newsPublisher:    newsPublisher,
	}
}

func (h *analysisHandler) HandleGetAnalysesByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	periods, err := h.analysis.GetAvailablePeriods(r.Context(), ticker)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get available periods")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"data": periods})
}

func (h *analysisHandler) HandleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		respondWithError(w, http.StatusBadRequest, "period query parameter is required (3, 6, 9, 12)")
		return
	}

	if _, ok := entity.MonthsToPeriod[period]; !ok {
		respondWithError(w, http.StatusBadRequest, "invalid period (allowed: 3, 6, 9, 12)")
		return
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		respondWithError(w, http.StatusBadRequest, "year query parameter is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid year parameter")
		return
	}

	periodInt, _ := strconv.Atoi(period)

	analysis, err := h.analysis.GetAnalysis(r.Context(), ticker, year, periodInt)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "analysis not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get analysis")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"data": analysis})
}

func (h *analysisHandler) HandleGetReportResults(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		respondWithError(w, http.StatusBadRequest, "period query parameter is required (3, 6, 9, 12)")
		return
	}

	if _, ok := entity.MonthsToPeriod[period]; !ok {
		respondWithError(w, http.StatusBadRequest, "invalid period (allowed: 3, 6, 9, 12)")
		return
	}

	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		respondWithError(w, http.StatusBadRequest, "year query parameter is required")
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid year parameter")
		return
	}

	periodInt, _ := strconv.Atoi(period)

	results, err := h.reportResults.GetReportResults(r.Context(), ticker, year, periodInt)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "report results not found")
			return
		}
		slog.Error("GetReportResults failed", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to get report results")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"data": results})
}

func (h *analysisHandler) HandleGetBusinessResearch(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	research, err := h.businessResearch.GetBusinessResearch(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "business research not found")
			return
		}
		slog.Error("GetBusinessResearch failed", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to get business research")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"data": research})
}

func (h *analysisHandler) HandleGetLatestReportResults(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	results, err := h.reportResults.GetLatestReportResults(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			respondWithError(w, http.StatusNotFound, "report results not found")
			return
		}
		slog.Error("GetLatestReportResults failed", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to get report results")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"data": results})
}

func (h *analysisHandler) HandleGetNews(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	news, err := h.news.GetFreshNews(r.Context(), ticker, 72*time.Hour)
	if err != nil {
		slog.Error("GetFreshNews failed", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to get news")
		return
	}

	if news == nil {
		respondWithError(w, http.StatusNotFound, "news not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"data": news})
}

func (h *analysisHandler) HandleTriggerNews(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	shouldContinue := false
	task := entity.Task{
		Id:             ticker,
		Ticker:         ticker,
		Type:           entity.NewsResearch,
		ShouldContinue: &shouldContinue,
	}

	payload, err := json.Marshal(task)
	if err != nil {
		slog.Error("HandleTriggerNews: marshal task", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	if err := h.newsPublisher.PublishMessage(r.Context(), payload); err != nil {
		slog.Error("HandleTriggerNews: publish message", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to trigger news fetch")
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{"status": "triggered"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
