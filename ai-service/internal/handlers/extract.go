package handlers

import (
	"ai-service/internal/application"
	"ai-service/internal/domain"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ExtractorHandler struct {
	service *application.ExtractorService
}

func NewExtractorHandler(service *application.ExtractorService) *ExtractorHandler {
	return &ExtractorHandler{service: service}
}

func (h *ExtractorHandler) HandleExtract(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")
	if ticker == "" {
		respondWithError(w, http.StatusBadRequest, "ticker query parameter is required")
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		respondWithError(w, http.StatusBadRequest, "period query parameter is required (3, 6, 9, 12)")
		return
	}

	if _, ok := domain.MonthsToPeriod[period]; !ok {
		respondWithError(w, http.StatusBadRequest, "invalid period (allowed: 3, 6, 9, 12)")
		return
	}

	var year int
	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		var err error
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid year parameter")
			return
		}
	}

	force := r.URL.Query().Get("force") == "true"

	rawData, err := h.service.Extract(r.Context(), ticker, period, year, force)
	if err != nil {
		log.Printf("Extraction failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("extraction failed: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rawData); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to encode response")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
