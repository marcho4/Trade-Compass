package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
			RespondWithError(w, r, 400, "invalid sector_id", parseErr)
			return
		}

		sector := domain.Sector(sectorID)
		if !sector.IsValid() {
			RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
			return
		}

		newsList, err = h.repo.GetBySector(r.Context(), sectorID)
	} else {
		RespondWithError(w, r, 400, "ticker or sector query parameter is required", nil)
		return
	}

	if err != nil {
		RespondWithError(w, r, 500, "failed to load news", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newsList); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *NewsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	news, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load news", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(news); err != nil {
		RespondWithError(w, r, 500, "failed to encode response", err)
		return
	}
}

func (h *NewsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	var news domain.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if news.SectorID != nil {
		sector := domain.Sector(*news.SectorID)
		if !sector.IsValid() {
			RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
			return
		}
	}

	if err := h.repo.Create(r.Context(), &news); err != nil {
		RespondWithError(w, r, 500, "failed to create news", err)
		return
	}

	RespondWithSuccess(w, 201, news, "News successfully created")
}

func (h *NewsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	var news domain.News
	if err := json.NewDecoder(r.Body).Decode(&news); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if news.SectorID != nil {
		sector := domain.Sector(*news.SectorID)
		if !sector.IsValid() {
			RespondWithError(w, r, 400, "invalid sector_id (allowed values from 1 to 19)", nil)
			return
		}
	}

	if err := h.repo.Update(r.Context(), id, &news); err != nil {
		RespondWithError(w, r, 500, "failed to update news", err)
		return
	}

	RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "News successfully updated")
}

func (h *NewsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		RespondWithError(w, r, 500, "failed to delete news", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "News successfully deleted")
}
