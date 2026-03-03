package routers

import (
	"encoding/json"
	"financial_data/internal/application/middleware"
	"financial_data/internal/application/response"
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type SectorHandler struct {
	repo SectorRepository
}

func NewSectorHandler(repo SectorRepository) *SectorHandler {
	return &SectorHandler{repo: repo}
}

func RegisterSectorRoutes(r chi.Router, repo SectorRepository, m *middleware.MiddlewareConfig) {
	handler := NewSectorHandler(repo)

	r.Get("/sectors", handler.HandleGetAll)
	r.Get("/sectors/{id}", handler.HandleGetByID)

	r.Group(func(protected chi.Router) {
		protected.Use(m.AuthMiddleware)
		protected.Post("/sectors", handler.HandleCreate)
		protected.Put("/sectors/{id}", handler.HandleUpdate)
		protected.Delete("/sectors/{id}", handler.HandleDelete)
	})
}

func (h *SectorHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	sector, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load sector", err)
		return
	}

	response.RespondWithSuccess(w, 200, sector, "Successfully retrieved sector by id")
}

func (h *SectorHandler) HandleGetAll(w http.ResponseWriter, r *http.Request) {
	sectors, err := h.repo.GetAll(r.Context())
	if err != nil {
		response.RespondWithError(w, r, 500, "failed to load sectors", err)
		return
	}

	response.RespondWithSuccess(w, 200, sectors, "Successfully retrieved all sectors")
}

func (h *SectorHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var sector domain.SectorModel
	if err := json.NewDecoder(r.Body).Decode(&sector); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if err := h.repo.Create(r.Context(), &sector); err != nil {
		response.RespondWithError(w, r, 500, "failed to create sector", err)
		return
	}

	response.RespondWithSuccess(w, 201, sector, "Sector successfully created")
}

func (h *SectorHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	var sector domain.SectorModel
	if err := json.NewDecoder(r.Body).Decode(&sector); err != nil {
		response.RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if err := h.repo.Update(r.Context(), id, &sector); err != nil {
		response.RespondWithError(w, r, 500, "failed to update sector", err)
		return
	}

	response.RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "Sector successfully updated")
}

func (h *SectorHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		response.RespondWithError(w, r, 400, "id is required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		response.RespondWithError(w, r, 500, "failed to delete sector", err)
		return
	}

	response.RespondWithSuccess(w, 204, nil, "Sector successfully deleted")
}
