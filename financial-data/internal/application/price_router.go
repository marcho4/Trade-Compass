package application

import (
	"financial_data/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PriceHandler struct {
	priceProvider domain.MarketService
}

func NewPriceHandler(priceProvider domain.MarketService) *PriceHandler {
	return &PriceHandler{priceProvider: priceProvider}
}

func RegisterPriceRoutes(r chi.Router, priceProvider domain.MarketService) {
	handler := NewPriceHandler(priceProvider)
	r.Get("/price", handler.HandleGetPriceByTicker)
	r.Get("/price/latest", handler.HandleGetLatestPrice)
	r.Get("/market-cap", handler.HandleGetMarketCap)
	r.Get("/stock-info", handler.HandleGetStockInfo)
}

func (h *PriceHandler) HandleGetPriceByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	days := r.URL.Query().Get("days")
	interval := r.URL.Query().Get("interval")
	if ticker == "" || days == "" || interval == "" {
		RespondWithError(w, r, http.StatusBadRequest, "ticker, days and interval are required in query params", nil)
		return
	}

	if !isValidInterval(interval) {
		RespondWithError(w, r, http.StatusBadRequest, "invalid interval in query params. Must be 60, 24 or 7", nil)
		return
	}

	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, "invalid interval in query params", err)
		return
	}
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		RespondWithError(w, r, http.StatusBadRequest, "invalid days in query params", err)
		return
	}
	price, err := h.priceProvider.GetStockPrice(ticker, daysInt, domain.Period(intervalInt))
	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, "failed to get price", err)
		return
	}

	RespondWithSuccess(w, http.StatusOK, price, "Success")
}

func (h *PriceHandler) HandleGetLatestPrice(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		RespondWithError(w, r, http.StatusBadRequest, "ticker is required", nil)
		return
	}

	price, err := h.priceProvider.GetStockPrice(ticker, 5, domain.Period(60))
	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, "failed to get latest price", err)
		return
	}

	RespondWithSuccess(w, http.StatusOK, price[len(price)-1].Close, "Success")
}

func (h *PriceHandler) HandleGetMarketCap(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		RespondWithError(w, r, http.StatusBadRequest, "ticker is required", nil)
		return
	}

	marketCap, err := h.priceProvider.GetMarketCap(ticker)
	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, "failed to get market cap", err)
		return
	}

	RespondWithSuccess(w, http.StatusOK, marketCap, "Successfully got market cap")
}

func (h *PriceHandler) HandleGetStockInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		RespondWithError(w, r, http.StatusBadRequest, "provide ticker in the url params", nil)
		return
	}

	stockInfo, err := h.priceProvider.GetStockInfo(ticker)
	if err != nil {
		RespondWithError(w, r, http.StatusInternalServerError, "Error happened while retrieving data", err)
		return
	}

	RespondWithSuccess(w, http.StatusOK, stockInfo, "Successfully got stock info")
}

func isValidInterval(i string) bool {
	switch i {
	case "60", "24", "7":
		return true
	default:
		return false
	}
}
