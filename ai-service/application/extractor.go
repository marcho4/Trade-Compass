package application

import (
	"ai-service/domain"
	"ai-service/infrastructure/financialdata"
	"ai-service/infrastructure/gemini"
	"ai-service/infrastructure/parser"
	"ai-service/infrastructure/s3"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type ExtractorService struct {
	gemini        *gemini.Client
	s3            *s3.Client
	parser        *parser.Client
	financialData *financialdata.Client
}

func NewExtractorService(
	geminiClient *gemini.Client,
	s3Client *s3.Client,
	parserClient *parser.Client,
	fdClient *financialdata.Client,
) *ExtractorService {
	return &ExtractorService{
		gemini:        geminiClient,
		s3:            s3Client,
		parser:        parserClient,
		financialData: fdClient,
	}
}

func (s *ExtractorService) Extract(ctx context.Context, ticker, period string, year int, force bool) (*domain.RawData, error) {
	mappedPeriod, ok := domain.MonthsToPeriod[period]
	if !ok {
		return nil, fmt.Errorf("invalid period: %s (expected: 3, 6, 9, 12)", period)
	}

	if year == 0 {
		latestYear, err := s.parser.GetLatestReportYear(ctx, ticker, period)
		if err != nil {
			return nil, fmt.Errorf("failed to determine latest year: %w", err)
		}
		year = latestYear
	}

	if !force {
		existing, err := s.financialData.GetDraft(ctx, ticker, year, mappedPeriod)
		if err != nil {
			log.Printf("Warning: failed to check existing draft: %v", err)
		}
		if existing != nil {
			log.Printf("Returning cached draft for %s %d/%s", ticker, year, period)
			return existing, nil
		}
	}

	s3Path, err := s.parser.GetReportS3Path(ctx, ticker, period, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get report S3 path: %w", err)
	}

	pdfBytes, err := s.s3.DownloadPDF(ctx, s3Path)
	if err != nil {
		return nil, fmt.Errorf("failed to download PDF: %w", err)
	}

	rawData, err := s.gemini.ExtractFromPDF(ctx, pdfBytes, ticker, year, period)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from PDF: %w", err)
	}

	if err := s.financialData.SaveDraft(ctx, rawData); err != nil {
		log.Printf("Warning: failed to save draft, trying update: %v", err)
		if err := s.financialData.UpdateDraft(ctx, rawData); err != nil {
			log.Printf("Warning: failed to update draft: %v", err)
		}
	}

	return rawData, nil
}

type ExtractorHandler struct {
	service *ExtractorService
}

func NewExtractorHandler(service *ExtractorService) *ExtractorHandler {
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
