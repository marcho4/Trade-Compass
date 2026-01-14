package application

import (
	"context"
	"encoding/json"
	"net/http"

	"financial_data/domain"

	"github.com/go-chi/chi"
)

type RatiosRepository interface {
	GetByTicker(ctx context.Context, ticker string) (*domain.Ratios, error)
}

type RatiosHandler struct {
	repo RatiosRepository
}

func NewRatiosHandler(repo RatiosRepository) *RatiosHandler {
	return &RatiosHandler{repo: repo}
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
