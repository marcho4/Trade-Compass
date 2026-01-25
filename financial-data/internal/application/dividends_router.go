package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type DividendsHandler struct {
	repo DividendsRepository
}

func NewDividendsHandler(repo DividendsRepository) *DividendsHandler {
	return &DividendsHandler{repo: repo}
}

func RegisterDividendsRoutes(r chi.Router, repo DividendsRepository) {
	handler := NewDividendsHandler(repo)

	r.Get("/dividends/{ticker}", handler.HandleGetByTicker)
	r.Get("/dividends/{ticker}/{id}", handler.HandleGetByID)
	r.Post("/dividends/{ticker}", handler.HandleCreate)
	r.Put("/dividends/{ticker}/{id}", handler.HandleUpdate)
	r.Delete("/dividends/{ticker}/{id}", handler.HandleDelete)
}

func (h *DividendsHandler) HandleGetByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	dividends, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load dividends: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dividends); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DividendsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		http.Error(w, "ticker and id are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	dividend, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load dividend: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dividend); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *DividendsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		http.Error(w, "ticker is required", http.StatusBadRequest)
		return
	}

	var dividend domain.Dividends
	if err := json.NewDecoder(r.Body).Decode(&dividend); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dividend.Ticker = ticker

	if err := h.repo.Create(r.Context(), &dividend); err != nil {
		http.Error(w, fmt.Sprintf("failed to create dividend: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dividend)
}

func (h *DividendsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		http.Error(w, "ticker and id are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var dividend domain.Dividends
	if err := json.NewDecoder(r.Body).Decode(&dividend); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dividend.Ticker = ticker

	if err := h.repo.Update(r.Context(), id, &dividend); err != nil {
		http.Error(w, fmt.Sprintf("failed to update dividend: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *DividendsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		http.Error(w, "ticker and id are required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete dividend: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
