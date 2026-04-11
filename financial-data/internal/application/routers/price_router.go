package routers

import (
	"financial_data/internal/application/middleware"
	"financial_data/internal/application/response"
	"financial_data/internal/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type PriceHandler struct {
	priceProvider domain.MarketService
}

func NewPriceHandler(priceProvider domain.MarketService) *PriceHandler {
	return &PriceHandler{priceProvider: priceProvider}
}

func RegisterPriceRoutes(r chi.Router, priceProvider domain.MarketService, m *middleware.MiddlewareConfig) {
	handler := NewPriceHandler(priceProvider)
	r.Get("/price", handler.HandleGetPriceByTicker)
	r.Get("/price/latest", handler.HandleGetLatestPrice)
	r.Get("/price/at", handler.HandleGetPriceAt)
	r.Get("/market-cap", handler.HandleGetMarketCap)
	r.Get("/stock-info", handler.HandleGetStockInfo)
}

func (h *PriceHandler) HandleGetPriceByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	days := r.URL.Query().Get("days")
	interval := r.URL.Query().Get("interval")
	if ticker == "" || days == "" || interval == "" {
		response.RespondWithError(w, r, http.StatusBadRequest, "ticker, days and interval are required in query params", nil)
		return
	}

	if !isValidInterval(interval) {
		response.RespondWithError(w, r, http.StatusBadRequest, "invalid interval in query params. Must be 60, 24 or 7", nil)
		return
	}

	intervalInt, err := strconv.Atoi(interval)
	if err != nil {
		response.RespondWithError(w, r, http.StatusBadRequest, "invalid interval in query params", err)
		return
	}
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		response.RespondWithError(w, r, http.StatusBadRequest, "invalid days in query params", err)
		return
	}
	price, err := h.priceProvider.GetStockPrice(ticker, daysInt, domain.Period(intervalInt))
	if err != nil {
		response.RespondWithError(w, r, http.StatusInternalServerError, "failed to get price", err)
		return
	}

	response.RespondWithSuccess(w, http.StatusOK, price, "Success")
}

func (h *PriceHandler) HandleGetLatestPrice(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		response.RespondWithError(w, r, http.StatusBadRequest, "ticker is required", nil)
		return
	}

	price, err := h.priceProvider.GetStockPrice(ticker, 5, domain.Period(60))
	if err != nil {
		response.RespondWithError(w, r, http.StatusInternalServerError, "failed to get latest price", err)
		return
	}

	response.RespondWithSuccess(w, http.StatusOK, price[len(price)-1].Close, "Success")
}

func (h *PriceHandler) HandleGetPriceAt(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	dateStr := r.URL.Query().Get("date")
	if ticker == "" || dateStr == "" {
		response.RespondWithError(w, r, http.StatusBadRequest, "ticker and date are required", nil)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.RespondWithError(w, r, http.StatusBadRequest, "date must be in YYYY-MM-DD format", err)
		return
	}

	price, err := h.priceProvider.GetPriceAt(ticker, date)
	if err != nil {
		response.RespondWithError(w, r, http.StatusInternalServerError, "failed to get price at date", err)
		return
	}

	response.RespondWithSuccess(w, http.StatusOK, price, "Success")
}

func (h *PriceHandler) HandleGetMarketCap(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		response.RespondWithError(w, r, http.StatusBadRequest, "ticker is required", nil)
		return
	}

	marketCap, err := h.priceProvider.GetMarketCap(ticker)
	if err != nil {
		response.RespondWithError(w, r, http.StatusInternalServerError, "failed to get market cap", err)
		return
	}

	response.RespondWithSuccess(w, http.StatusOK, marketCap, "Successfully got market cap")
}

func (h *PriceHandler) HandleGetStockInfo(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		response.RespondWithError(w, r, http.StatusBadRequest, "provide ticker in the url params", nil)
		return
	}

	stockInfo, err := h.priceProvider.GetStockInfo(ticker)
	if err != nil {
		response.RespondWithError(w, r, http.StatusInternalServerError, "Error happened while retrieving data", err)
		return
	}

	response.RespondWithSuccess(w, http.StatusOK, stockInfo, "Successfully got stock info")
}

func isValidInterval(i string) bool {
	switch i {
	case "60", "24", "7":
		return true
	default:
		return false
	}
}
