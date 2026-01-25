package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type NewsHandler struct {
	repo NewsRepository
}

func NewNewsHandler(repo NewsRepository) *NewsHandler {
	return &NewsHandler{repo: repo}
}

func RegisterNewsRoutes(r chi.Router, repo NewsRepository) {
	handler := NewNewsHandler(repo)

	r.Get("/news", handler.HandleGetNews)
	r.Get("/news/{id}", handler.HandleGetByID)
	r.Post("/news", handler.HandleCreate)
	r.Put("/news/{id}", handler.HandleUpdate)
	r.Delete("/news/{id}", handler.HandleDelete)
}

func (h *NewsHandler) HandleGetNews(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	sectorIDStr := r.URL.Query().Get("sector")

	var newsList []domain.News
	var err error

	if ticker != "" {
		newsList, err = h.repo.GetByTicker(r.Context(), ticker)
	} else if sectorIDStr != "" {
		sectorID, parseErr := strconv.Atoi(sectorIDStr)
		if parseErr != nil {
			http.Error(w, "invalid sector_id", http.StatusBadRequest)
			return
		}

		sector := domain.Sector(sectorID)
		if !sector.IsValid() {
			http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
			return
		}

		newsList, err = h.repo.GetBySector(r.Context(), sectorID)
	} else {
		http.Error(w, "ticker or sector query parameter is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load news: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newsList); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *NewsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
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

	news, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load news: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(news); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *NewsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var news domain.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if news.SectorID != nil {
		sector := domain.Sector(*news.SectorID)
		if !sector.IsValid() {
			http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
			return
		}
	}

	if err := h.repo.Create(r.Context(), &news); err != nil {
		http.Error(w, fmt.Sprintf("failed to create news: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}

func (h *NewsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
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

	var news domain.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if news.SectorID != nil {
		sector := domain.Sector(*news.SectorID)
		if !sector.IsValid() {
			http.Error(w, "invalid sector_id (allowed values from 1 to 19)", http.StatusBadRequest)
			return
		}
	}

	if err := h.repo.Update(r.Context(), id, &news); err != nil {
		http.Error(w, fmt.Sprintf("failed to update news: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *NewsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, fmt.Sprintf("failed to delete news: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
