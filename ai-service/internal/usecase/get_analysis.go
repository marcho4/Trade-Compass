package usecase

import (
	"context"

	"ai-service/internal/domain/entity"
)

type GetAnalysisUsecase struct {
	repo AnalysisRepository
}

func NewGetAnalysisUsecase(repo AnalysisRepository) *GetAnalysisUsecase {
	return &GetAnalysisUsecase{repo: repo}
}

func (u *GetAnalysisUsecase) GetAnalysis(ctx context.Context, ticker string, year, period int) (string, error) {
	return u.repo.GetAnalysis(ctx, ticker, year, period)
}

func (u *GetAnalysisUsecase) GetAvailablePeriods(ctx context.Context, ticker string) ([]entity.AvailablePeriod, error) {
	return u.repo.GetAvailablePeriods(ctx, ticker)
}
