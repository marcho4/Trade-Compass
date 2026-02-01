package application

import (
	"encoding/json"
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

func RegisterRatiosRoutes(r chi.Router, repo RatiosRepository) {
	handler := NewRatiosHandler(repo)

	r.Get("/ratios/sector/{sector_id}", handler.HandleGetRatiosBySector)
	r.Get("/ratios/{ticker}", handler.HandleGetRatiosByTicker)
	r.Post("/ratios/{ticker}", handler.HandleCreateRatios)
	r.Put("/ratios/{ticker}", handler.HandleUpdateRatios)
	r.Delete("/ratios/{ticker}", handler.HandleDeleteRatios)
}

func (h *RatiosHandler) HandleGetRatiosByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	ratios, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load ratios", err)
		return
	}

	RespondWithSuccess(w, 200, ratios, "Successfully retrieved ratios by ticker")
}

func (h *RatiosHandler) HandleGetRatiosBySector(w http.ResponseWriter, r *http.Request) {
	sector_id := chi.URLParam(r, "sector_id")
	if sector_id == "" {
		RespondWithError(w, r, 400, "sector_id is required", nil)
		return
	}

	parsedSector, err := strconv.Atoi(sector_id)
	if err != nil {
		RespondWithError(w, r, 400, "sector_id must be int", err)
		return
	}

	sector := domain.Sector(parsedSector)
	if !sector.IsValid() {
		RespondWithError(w, r, 400, "Sector is not valid (allowed values from 1 to 19)", nil)
		return
	}

	ratios, err := h.repo.GetBySector(r.Context(), domain.Sector(parsedSector))
	if err != nil {
		RespondWithError(w, r, 500, "failed to load ratios", err)
		return
	}

	RespondWithSuccess(w, 200, ratios, "Successfully retrieved ratios by sector")
}

func (h *RatiosHandler) HandleCreateRatios(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var requestBody struct {
		SectorID int           `json:"sector_id"`
		Ratios   domain.Ratios `json:"ratios"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	sector := domain.Sector(requestBody.SectorID)
	if !sector.IsValid() {
		RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
		return
	}

	if err := h.repo.Create(r.Context(), ticker, sector, &requestBody.Ratios); err != nil {
		RespondWithError(w, r, 500, "failed to create ratios", err)
		return
	}

	RespondWithSuccess(w, 201, map[string]string{"status": "created"}, "Ratios successfully created")
}

func (h *RatiosHandler) HandleUpdateRatios(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var ratios domain.Ratios
	if err := json.NewDecoder(r.Body).Decode(&ratios); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if err := h.repo.Update(r.Context(), ticker, &ratios); err != nil {
		RespondWithError(w, r, 500, "failed to update ratios", err)
		return
	}

	RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "Ratios successfully updated")
}

func (h *RatiosHandler) HandleDeleteRatios(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker); err != nil {
		RespondWithError(w, r, 500, "failed to delete ratios", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "Ratios successfully deleted")
}
