package routers

import (
	"encoding/json"
	"errors"
	"financial_data/internal/application/middleware"
	"financial_data/internal/application/response"
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RatiosHandler struct {
	repo RatiosRepository
}

func NewRatiosHandler(repo RatiosRepository) *RatiosHandler {
	return &RatiosHandler{repo: repo}
}

func RegisterRatiosRoutes(r chi.Router, repo RatiosRepository, m *middleware.MiddlewareConfig) {
	handler := NewRatiosHandler(repo)

	r.Get("/ratios/sector/{sector_id}", handler.HandleGetBySector)
	r.Get("/ratios/{ticker}", handler.HandleGetByPeriod)
	r.Get("/ratios/{ticker}/latest", handler.HandleGetLatest)
	r.Get("/ratios/{ticker}/history", handler.HandleGetHistory)

	r.Group(func(protected chi.Router) {
		protected.Use(m.AuthMiddleware)

		protected.Post("/ratios/{ticker}", handler.HandleCreate)
		protected.Put("/ratios/{ticker}", handler.HandleUpdate)
		protected.Delete("/ratios/{ticker}", handler.HandleDelete)
	})
}

func (h *RatiosHandler) HandleGetByPeriod(w http.ResponseWriter, r *http.Request) {
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

	ratios, err := h.repo.GetByTickerAndPeriod(r.Context(), ticker, year, period)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondWithError(w, r, 404, "ratios not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to load ratios", err)
		return
	}

	response.RespondWithSuccess(w, 200, ratios, "")
}

func (h *RatiosHandler) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	ratios, err := h.repo.GetLatestByTicker(r.Context(), ticker)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.RespondWithError(w, r, 404, "ratios not found", err)
			return
		}
		response.RespondWithError(w, r, 500, "failed to load latest ratios", err)
		return
	}

	response.RespondWithSuccess(w, 200, ratios, "")
}

func (h *RatiosHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	history, err := h.repo.GetHistoryByTicker(r.Context(), ticker)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load ratios history", err)
		return
	}

	if history == nil {
		history = []domain.Ratios{}
	}

	response.RespondWithSuccess(w, 200, history, "")
}

func (h *RatiosHandler) HandleGetBySector(w http.ResponseWriter, r *http.Request) {
	sectorID := chi.URLParam(r, "sector_id")
	if sectorID == "" {
		response.RespondWithError(w, r, 400, "sector_id is required", nil)
		return
	}

	parsed, err := strconv.Atoi(sectorID)
	if err != nil {
		response.RespondWithError(w, r, 400, "sector_id must be int", err)
		return
	}

	sector := domain.Sector(parsed)
	if !sector.IsValid() {
		response.RespondWithError(w, r, 400, "Sector is not valid (allowed values from 1 to 19)", nil)
		return
	}

	ratios, err := h.repo.GetBySector(r.Context(), sector)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load ratios", err)
		return
	}

	response.RespondWithSuccess(w, 200, ratios, "")
}

func (h *RatiosHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var requestBody struct {
		SectorID int           `json:"sector_id"`
		Ratios   domain.Ratios `json:"ratios"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(requestBody.SectorID)
	if !sector.IsValid() {
		response.RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	requestBody.Ratios.Ticker = ticker

	if err := h.repo.Create(r.Context(), sector, &requestBody.Ratios); err != nil {
		response.RespondWithError(w, r, 500, "failed to create ratios", err)
		return
	}

	response.RespondWithSuccess(w, 201, map[string]string{"status": "created"}, "Ratios successfully created")
}

func (h *RatiosHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		response.RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var ratios domain.Ratios
	if err := json.NewDecoder(r.Body).Decode(&ratios); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	ratios.Ticker = ticker

	if err := h.repo.Update(r.Context(), &ratios); err != nil {
		response.RespondWithError(w, r, 500, "failed to update ratios", err)
		return
	}

	response.RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "Ratios successfully updated")
}

func (h *RatiosHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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
		response.RespondWithError(w, r, 500, "failed to delete ratios", err)
		return
	}

	response.RespondWithSuccess(w, 204, nil, "Ratios successfully deleted")
}
