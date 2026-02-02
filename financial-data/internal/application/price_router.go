package application

import (
	"financial_data/internal/domain"
	"financial_data/internal/infrastructure"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PriceHandler struct {
	priceProvider *infrastructure.MoexDataProvider
}

func NewPriceHandler(priceProvider *infrastructure.MoexDataProvider) *PriceHandler {
	return &PriceHandler{priceProvider: priceProvider}
}

func RegisterPriceRoutes(r chi.Router, priceProvider *infrastructure.MoexDataProvider) {
	handler := NewPriceHandler(priceProvider)
	r.Get("/price", handler.HandleGetPriceByTicker)
	r.Get("/stock-info", handler.HandleGetStockInfo)
}

func (h *PriceHandler) HandleGetPriceByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	days := r.URL.Query().Get("days")
	interval := r.URL.Query().Get("interval")
	if ticker == "" || days == "" || interval == "" {
		RespondWithError(w, r, 400, "ticker, days and interval are required in query params", nil)
		return
	}
	if interval != "60" && interval != "24" && interval != "7" {
		RespondWithError(w, r, 400, "invalid interval in query params. Must be 60, 24 or 7", nil)
		return
	}
	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		RespondWithError(w, r, 400, "invalid interval in query params", err)
		return
	}
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		RespondWithError(w, r, 400, "invalid days in query params", err)
		return
	}
	price, err := h.priceProvider.GetStockPrice(ticker, daysInt, domain.Period(intervalInt))
	if err != nil {
		RespondWithError(w, r, 500, "failed to get price", err)
		return
	}

	RespondWithSuccess(w, 200, price, "Success")
}

func (h *PriceHandler) HandleGetStockInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		RespondWithError(w, r, 400, "provide ticker in the url params", nil)
	}
	stockInfo, err := h.priceProvider.GetStockInfo(ticker)
	if err != nil {
		RespondWithError(w, r, 500, "Error happened while retrieving data", err)
	}
	RespondWithSuccess(w, 200, stockInfo, "Successfully got stock info")
}
