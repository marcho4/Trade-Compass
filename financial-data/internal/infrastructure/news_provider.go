package infrastructure

import "financial_data/internal/domain"

type NewsProvider struct{}

func NewNewsProvider() *NewsProvider {
	return &NewsProvider{}
}

func GetTickerNews(ticker string) ([]domain.News, error) {
	return nil, nil
}

func GetSectorNews(sector domain.Sector) ([]domain.News, error) {
	return nil, nil
}
