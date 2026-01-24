package application

import (
	"encoding/json"
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type PriceHandler struct {
	priceProvider *infrastructure.MoexPriceProvider
}

func NewPriceHandler(priceProvider *infrastructure.MoexPriceProvider) *PriceHandler {
	return &PriceHandler{priceProvider: priceProvider}
}

func RegisterPriceRoutes(r chi.Router, priceProvider *infrastructure.MoexPriceProvider) {
	handler := NewPriceHandler(priceProvider)
	r.Get("/price", handler.HandleGetPriceByTicker)
}

func (h *PriceHandler) HandleGetPriceByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	days := r.URL.Query().Get("days")
	interval := r.URL.Query().Get("interval")
	if ticker == "" || days == "" || interval == "" {
		http.Error(w, "ticker, days and interval are required in query params", http.StatusBadRequest)
		return
	}
	if interval != "60" && interval != "24" && interval != "7" {
		http.Error(w, "invalid interval in query params. Must be 60, 24 or 7", http.StatusBadRequest)
		return
	}
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		http.Error(w, "invalid interval in query params", http.StatusBadRequest)
		return
	}
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		http.Error(w, "invalid days in query params", http.StatusBadRequest)
		return
	}
	price, err := h.priceProvider.GetStockPrice(ticker, daysInt, domain.Period(intervalInt))
	if err != nil {
		http.Error(w, "failed to get price", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(price); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
