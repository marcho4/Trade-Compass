package handlers

import (
	"ai-service/internal/domain"
	"ai-service/internal/infrastructure/postgres"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type AnalysisHandler struct {
	db *postgres.DBRepo
}

func NewAnalysisHandler(db *postgres.DBRepo) *AnalysisHandler {
	return &AnalysisHandler{db: db}
}

func (h *AnalysisHandler) HandleGetAnalysesByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	periods, err := h.db.GetAvailablePeriods(r.Context(), ticker)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get available periods")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": periods})
}

func (h *AnalysisHandler) HandleGetAnalysis(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := domain.MonthsToPeriod[period]; !ok {
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

	analysis, err := h.db.GetAnalysis(r.Context(), ticker, year, periodInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "analysis not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get analysis")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": analysis})
}

func (h *AnalysisHandler) HandleGetReportResults(w http.ResponseWriter, r *http.Request) {
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

	if _, ok := domain.MonthsToPeriod[period]; !ok {
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

	results, err := h.db.GetReportResults(r.Context(), ticker, year, periodInt)
	if err != nil {
		slog.Error("Error", slog.Any("err", err))
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "report results not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "failed to get report results")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": results})
}

func (h *AnalysisHandler) HandleGetLatestReportResults(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	results, err := h.db.GetLatestReportResults(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "report results not found")
			return
		}
		slog.Error("GetLatestReportResults failed", slog.String("ticker", ticker), slog.Any("error", err))
		respondWithError(w, http.StatusInternalServerError, "failed to get report results")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"data": results})
}
