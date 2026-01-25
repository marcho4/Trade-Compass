package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
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
		http.Error(w, fmt.Sprintf("failed to load current CB rate: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rate); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *MacroHandler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		http.Error(w, "from and to query parameters are required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		http.Error(w, "invalid from date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		http.Error(w, "invalid to date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	rates, err := h.repo.GetHistory(r.Context(), from, to)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load CB rates history: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rates); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *MacroHandler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var request struct {
		Date string  `json:"date"`
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		http.Error(w, "invalid date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	rate := &domain.CBRate{
		Date: date,
		Rate: request.Rate,
	}

	if err := h.repo.Create(r.Context(), rate); err != nil {
		http.Error(w, fmt.Sprintf("failed to create CB rate: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "created"})
}

func (h *MacroHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date query parameter is required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	var request struct {
		Rate float64 `json:"rate"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), date, request.Rate); err != nil {
		http.Error(w, fmt.Sprintf("failed to update CB rate: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *MacroHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if ok, err := validateApiKey(r); !ok {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(w, "date query parameter is required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "invalid date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), date); err != nil {
		http.Error(w, fmt.Sprintf("failed to delete CB rate: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
