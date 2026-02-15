package application

import (
	"ai-service/internal/domain"
	"ai-service/internal/infrastructure/financialdata"
	"ai-service/internal/infrastructure/gemini"
	"ai-service/internal/infrastructure/parser"
	"ai-service/internal/infrastructure/s3"
	"context"
	"fmt"
	"log"
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
