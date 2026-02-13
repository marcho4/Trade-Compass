package domain

import "context"

type ContextCollector interface{
	GetSectorNews(ctx context.Context, sectorName string) (string, error)
	GetCompanyNews(ctx context.Context, ticker string) (string, error)
	GetCompanyPriceHistory(ctx context.Context, ticker string) (string, error)
	GetMacroeconomicData(ctx context.Context, ticker string) (string, error)
	GetResearchFramework(ctx context.Context) (string, error)
	GetUserLevel(ctx context.Context, userID string) (string, error)
	GetUserPreferences(ctx context.Context, userID string) (string, error)
}