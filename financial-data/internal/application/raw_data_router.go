package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		http.Error(w, "year and period query parameters are required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		http.Error(w, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", http.StatusBadRequest)
		return
	}

	rawData, err := h.repo.GetByTickerAndPeriod(r.Context(), ticker, year, period)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *RawDataHandler) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	rawData, err := h.repo.GetLatestByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load latest metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *RawDataHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	history, err := h.repo.GetHistoryByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load metrics history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *RawDataHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rawData.Ticker = ticker

	if !rawData.Period.IsValid() {
		http.Error(w, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &rawData); err != nil {
		http.Error(w, fmt.Sprintf("failed to create metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *RawDataHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		http.Error(w, "year and period query parameters are required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		http.Error(w, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", http.StatusBadRequest)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rawData.Ticker = ticker
	rawData.Year = year
	rawData.Period = period

	if err := h.repo.Update(r.Context(), &rawData); err != nil {
		http.Error(w, fmt.Sprintf("failed to update metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *RawDataHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		http.Error(w, "year and period query parameters are required", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "invalid year parameter", http.StatusBadRequest)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		http.Error(w, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker, year, period); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete metrics: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
