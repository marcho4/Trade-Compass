package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	dividends, err := h.repo.GetByTicker(r.Context(), ticker)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load dividends", err)
		return
	}

	RespondWithSuccess(w, 200, dividends, "Successfully retrieved dividends by ticker")
}

func (h *DividendsHandler) HandleGetByID(w http.ResponseWriter, r *http.Request) {
	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		RespondWithError(w, r, 400, "ticker and id are required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	dividend, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, r, 500, "failed to load dividend", err)
		return
	}

	RespondWithSuccess(w, 200, dividend, "Successfully retrieved dividend by id")
}

func (h *DividendsHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "ticker is required", nil)
		return
	}

	var dividend domain.Dividends
	if err := json.NewDecoder(r.Body).Decode(&dividend); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	dividend.Ticker = ticker

	if err := h.repo.Create(r.Context(), &dividend); err != nil {
		RespondWithError(w, r, 500, "failed to create dividend", err)
		return
	}

	RespondWithSuccess(w, 201, dividend, "Dividend successfully created")
}

func (h *DividendsHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		RespondWithError(w, r, 400, "ticker and id are required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	var dividend domain.Dividends
	if err := json.NewDecoder(r.Body).Decode(&dividend); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	dividend.Ticker = ticker

	if err := h.repo.Update(r.Context(), id, &dividend); err != nil {
		RespondWithError(w, r, 500, "failed to update dividend", err)
		return
	}

	RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "Dividend successfully updated")
}

func (h *DividendsHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	ticker := chi.URLParam(r, "ticker")
	idStr := chi.URLParam(r, "id")

	if ticker == "" || idStr == "" {
		RespondWithError(w, r, 400, "ticker and id are required", nil)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid id", err)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		RespondWithError(w, r, 500, "failed to delete dividend", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "Dividend successfully deleted")
}
