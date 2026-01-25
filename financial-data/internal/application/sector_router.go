package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type SectorHandler struct {
	repo SectorRepository
}

func NewSectorHandler(repo SectorRepository) *SectorHandler {
	return &SectorHandler{repo: repo}
}

func RegisterSectorRoutes(r chi.Router, repo SectorRepository) {
	handler := NewSectorHandler(repo)

	r.Get("/sectors", handler.HandleGetAll)
	r.Get("/sectors/{id}", handler.HandleGetByID)
	r.Post("/sectors", handler.HandleCreate)
	r.Put("/sectors/{id}", handler.HandleUpdate)
	r.Delete("/sectors/{id}", handler.HandleDelete)
}

func (h *SectorHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	sector, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load sector: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sector); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *SectorHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	sectors, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load sectors: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(sectors); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *SectorHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var sector domain.SectorModel
	if err := json.NewDecoder(r.Body).Decode(&sector); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &sector); err != nil {
		http.Error(w, fmt.Sprintf("failed to create sector: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sector)
}

func (h *SectorHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var sector domain.SectorModel
	if err := json.NewDecoder(r.Body).Decode(&sector); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), id, &sector); err != nil {
		http.Error(w, fmt.Sprintf("failed to update sector: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *SectorHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete sector: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
