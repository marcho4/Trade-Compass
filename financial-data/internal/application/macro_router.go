package application

import (
	"encoding/json"
	"errors"
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type MacroHandler struct {
	repo MacroDataRepository
}

func NewMacroHandler(repo MacroDataRepository) *MacroHandler {
	return &MacroHandler{repo: repo}
}

func RegisterMacroRoutes(r chi.Router, repo MacroDataRepository) {
	handler := NewMacroHandler(repo)

	r.Get("/macro/cb-rate/current", handler.HandleGetCurrent)
	r.Get("/macro/cb-rate/history", handler.HandleGetHistory)
	r.Post("/macro/cb-rate", handler.HandleCreate)
	r.Put("/macro/cb-rate", handler.HandleUpdate)
	r.Delete("/macro/cb-rate", handler.HandleDelete)
}

func (h *MacroHandler) HandleGetCurrent(w http.ResponseWriter, r *http.Request) {
	rate, err := h.repo.GetCurrent(r.Context())
	if err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to load current CB rate", err)
		return
	}

	RespondWithSuccess(w, 200, rate, "Successfully got current central bank rate")
}

func (h *MacroHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		RespondWithError(w, r, 400, "from and to query parameters are required (format: YYYY-MM-DD)", nil)
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid from date format (expected YYYY-MM-DD)", err)
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid to date format (expected YYYY-MM-DD)", err)
		return
	}

	rates, err := h.repo.GetHistory(r.Context(), from, to)
	if err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to load CB rates history", err)
		return
	}

	RespondWithSuccess(w, 200, rates, "Successfully retrieved CB rates history")
}

func (h *MacroHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	var request struct {
		Date string  `json:"date"`
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		RespondWithError(w, r, 400, "invalid date format (expected YYYY-MM-DD)", err)
		return
	}

	rate := &domain.CBRate{
		Date: date,
		Rate: request.Rate,
	}

	if err := h.repo.Create(r.Context(), rate); err != nil {
		RespondWithError(w, r, 500, "failed to create CB rate", err)
		return
	}

	RespondWithSuccess(w, 201, rate, "CB rate successfully created")
}

func (h *MacroHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		RespondWithError(w, r, 400, "date query parameter is required (format: YYYY-MM-DD)", nil)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid date format (expected YYYY-MM-DD)", err)
		return
	}

	var request struct {
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		RespondWithError(w, r, 400, "invalid request body", err)
		return
	}

	if err := h.repo.Update(r.Context(), date, request.Rate); err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to update CB rate", err)
		return
	}

	RespondWithSuccess(w, 200, map[string]string{"status": "updated"}, "CB rate successfully updated")
}

func (h *MacroHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		RespondWithError(w, r, 401, "unauthorized", err)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		RespondWithError(w, r, 400, "date query parameter is required (format: YYYY-MM-DD)", nil)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		RespondWithError(w, r, 400, "invalid date format (expected YYYY-MM-DD)", err)
		return
	}

	if err := h.repo.Delete(r.Context(), date); err != nil {
		var dbErr *infrastructure.DbError
		if errors.As(err, &dbErr) && dbErr.RowsAffected == 0 {
			RespondWithError(w, r, 404, dbErr.Message, err)
			return
		}
		RespondWithError(w, r, 500, "failed to delete CB rate", err)
		return
	}

	RespondWithSuccess(w, 204, nil, "CB rate successfully deleted")
}
