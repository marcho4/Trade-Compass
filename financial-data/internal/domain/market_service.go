package domain

import "time"

type MarketService interface {
	GetMarketCap(ticker string) (float64, error)
	GetStockInfo(ticker string) (*StockInfo, error)
	GetStockPrice(ticker string, daysBackwards int, interval Period) ([]Candle, error)
	GetPriceAt(ticker string, date time.Time) (float64, error)
}
