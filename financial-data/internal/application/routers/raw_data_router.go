package routers

import (
	"encoding/json"
	"errors"
	"financial_data/internal/application/middleware"
	"financial_data/internal/application/response"
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RawDataHandler struct {
	repo *infrastructure.RawDataRepository
}

func NewRawDataHandler(repo *infrastructure.RawDataRepository) *RawDataHandler {
	return &RawDataHandler{repo: repo}
}

func RegisterRawDataRoutes(r chi.Router, repo *infrastructure.RawDataRepository, m *middleware.MiddlewareConfig) {
	handler := NewRawDataHandler(repo)

	r.Get("/raw-data/{ticker}", handler.HandleGetByPeriod)
	r.Get("/raw-data/{ticker}/latest", handler.HandleGetLatest)
	r.Get("/raw-data/{ticker}/history", handler.HandleGetHistory)
	r.Get("/raw-data/{ticker}/drafts", handler.HandleGetDrafts)
	r.Get("/raw-data/{ticker}/draft", handler.HandleGetDraft)

	r.Group(func(protected chi.Router) {
		protected.Use(m.AuthMiddleware)

		protected.Put("/raw-data/{ticker}/confirm", handler.HandleConfirmDraft)
		protected.Post("/raw-data/{ticker}", handler.HandleCreate)
		protected.Put("/raw-data/{ticker}", handler.HandleUpdate)
		protected.Delete("/raw-data/{ticker}", handler.HandleDelete)
	})
}

func (h *RawDataHandler) HandleGetByPeriod(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		response.RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	rawData, err := h.repo.GetByTickerAndPeriod(r.Context(), ticker, year, period)
	if err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.Message == "metrics not found" {
			response.RespondWithError(w, r, 404, "metrics not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to load metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	rawData, err := h.repo.GetLatestByTicker(r.Context(), ticker)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load latest metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	history, err := h.repo.GetHistoryByTicker(r.Context(), ticker)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load metrics history", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleGetDrafts(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	drafts, err := h.repo.GetDraftsByTicker(r.Context(), ticker)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load drafts", err)
		return
	}

	if drafts == nil {
		drafts = []domain.RawData{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(drafts); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
	}
}

func (h *RawDataHandler) HandleGetDraft(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		response.RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	draft, err := h.repo.GetDraftByTickerAndPeriod(r.Context(), ticker, year, period)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load draft", err)
		return
	}

	if draft == nil {
		response.RespondWithError(w, r, 404, "draft not found", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(draft); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
	}
}

func (h *RawDataHandler) HandleConfirmDraft(w http.ResponseWriter, r *http.Request) {

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		response.RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	if err := h.repo.ConfirmDraft(r.Context(), ticker, year, period); err != nil {
		response.RespondWithError(w, r, 500, "failed to confirm draft", err)
		return
	}

	response.RespondWithSuccess(w, 200, map[string]string{"status": "confirmed"}, "Draft confirmed successfully")
}

func (h *RawDataHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	rawData.Ticker = ticker

	if !rawData.Period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	if err := h.repo.Create(r.Context(), &rawData); err != nil {
		response.RespondWithError(w, r, 500, "failed to create metrics", err)
		return
	}

	response.RespondWithSuccess(w, 201, map[string]string{"status": "created"}, "Metrics successfully created")
}

func (h *RawDataHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		response.RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	var rawData domain.RawData
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	rawData.Ticker = ticker
	rawData.Year = year
	rawData.Period = period

	if err := h.repo.Update(r.Context(), &rawData); err != nil {
		response.RespondWithError(w, r, 500, "failed to update metrics", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "updated"}); err != nil {
		response.RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *RawDataHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	yearStr := r.URL.Query().Get("year")
	periodStr := r.URL.Query().Get("period")

	if yearStr == "" || periodStr == "" {
		response.RespondWithError(w, r, 400, "year and period query parameters are required", nil)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid year parameter", err)
		return
	}

	period := domain.ReportPeriod(periodStr)
	if !period.IsValid() {
		response.RespondWithError(w, r, 400, "invalid period (allowed: Q1, Q2, Q3, Q4, YEAR)", nil)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker, year, period); err != nil {
		response.RespondWithError(w, r, 500, "failed to delete metrics", err)
		return
	}

	response.RespondWithSuccess(w, 204, nil, "Metrics successfully deleted")
}
