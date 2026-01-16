package infrastructure

import "financial_data/internal/domain"

type NewsProvider struct{}

func NewNewsProvider() *NewsProvider {
	return &NewsProvider{}
}

func GetTickerNews(ticker string) {}

func GetSectorNews(sector domain.Sector) {}
