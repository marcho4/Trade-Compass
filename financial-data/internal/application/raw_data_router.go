package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RawDataHandler struct {
	repo RawDataRepository
}

func NewRawDataHandler(repo RawDataRepository) *RawDataHandler {
	return &RawDataHandler{repo: repo}
}

func RegisterRawDataRoutes(r chi.Router, repo RawDataRepository) {
	handler := NewRawDataHandler(repo)

	r.Get("/raw-data/{ticker}", handler.HandleGetByPeriod)
	r.Get("/raw-data/{ticker}/latest", handler.HandleGetLatest)
	r.Get("/raw-data/{ticker}/history", handler.HandleGetHistory)
	r.Post("/raw-data/{ticker}", handler.HandleCreate)
	r.Put("/raw-data/{ticker}", handler.HandleUpdate)
	r.Delete("/raw-data/{ticker}", handler.HandleDelete)
}

func (h *RawDataHandler) HandleGetByPeriod(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	rawData, err := h.repo.GetByTickerAndPeriod(r.Context(), ticker, year, period)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	rawData, err := h.repo.GetLatestByTicker(r.Context(), ticker)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load latest metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	history, err := h.repo.GetHistoryByTicker(r.Context(), ticker)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load metrics history", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	rawData.Ticker = ticker

	if !rawData.Period.IsValid() {
		RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	if err := h.repo.Create(r.Context(), &rawData); err != nil {
		RespondWithError(w, r, 500, "failed to create metrics", err)
		return
	}

	RespondWithSuccess(w, 201, map[string]string{"status": "created"}, "Metrics successfully created")
}

func (h *RawDataHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	rawData.Ticker = ticker
	rawData.Year = year
	rawData.Period = period

	if err := h.repo.Update(r.Context(), &rawData); err != nil {
		RespondWithError(w, r, 500, "failed to update metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "updated"}); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker, year, period); err != nil {
		RespondWithError(w, r, 500, "failed to delete metrics", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "Metrics successfully deleted")
}
