package usecase

import (
	"context"

	"ai-service/internal/domain/entity"
)

type ReportResultsReader interface {
	GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error)
	GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error)
}

type GetReportResultsUsecase struct {
	repo ReportResultsReader
}

func NewGetReportResultsUsecase(repo ReportResultsReader) *GetReportResultsUsecase {
	return &GetReportResultsUsecase{repo: repo}
}

func (u *GetReportResultsUsecase) GetReportResults(ctx context.Context, ticker string, year, period int) (*entity.ReportResults, error) {
	return u.repo.GetReportResults(ctx, ticker, year, period)
}

func (u *GetReportResultsUsecase) GetLatestReportResults(ctx context.Context, ticker string) (*entity.ReportResults, error) {
	return u.repo.GetLatestReportResults(ctx, ticker)
}
