package application

import (
	"context"
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type RatiosRepository interface {
	GetByTicker(ctx context.Context, ticker string) (*domain.Ratios, error)
	GetBySector(ctx context.Context, sector domain.Sector) (*domain.Ratios, error)
	Create(ctx context.Context, ticker string, sector domain.Sector, ratios *domain.Ratios) error
	Update(ctx context.Context, ticker string, ratios *domain.Ratios) error
	Delete(ctx context.Context, ticker string) error
}

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
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	ratios, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, "failed to load ratios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ratios); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *RatiosHandler) HandleGetRatiosBySector(w http.ResponseWriter, r *http.Request) {
	sector_id := chi.URLParam(r, "sector_id")
	if sector_id == "" {
		http.Error(w, "sector_id is required", http.StatusBadRequest)
		return
	}

	parsedSector, err := strconv.Atoi(sector_id)
	if err != nil {
		http.Error(w, "sector_id must be int", http.StatusBadRequest)
		return
	}

	sector := domain.Sector(parsedSector)
	if !sector.IsValid() {
		http.Error(w, "Sector is not valid (allowed values from 1 to 19)", http.StatusBadRequest)
		return
	}

	ratios, err := h.repo.GetBySector(r.Context(), domain.Sector(parsedSector))
	if err != nil {
		http.Error(w, "failed to load ratios", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ratios); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *RatiosHandler) HandleCreateRatios(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		SectorID int           `json:"sector_id"`
		Ratios   domain.Ratios `json:"ratios"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	sector := domain.Sector(requestBody.SectorID)
	if !sector.IsValid() {
		http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), ticker, sector, &requestBody.Ratios); err != nil {
		http.Error(w, fmt.Sprintf("failed to create ratios: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *RatiosHandler) HandleUpdateRatios(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	var ratios domain.Ratios
	if err := json.NewDecoder(r.Body).Decode(&ratios); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), ticker, &ratios); err != nil {
		http.Error(w, fmt.Sprintf("failed to update ratios: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *RatiosHandler) HandleDeleteRatios(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), ticker); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete ratios: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
