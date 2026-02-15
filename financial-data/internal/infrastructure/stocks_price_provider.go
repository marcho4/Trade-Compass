package infrastructure

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MoexDataProvider struct {
	baseUrl string
	client  http.Client
}

func NewMoexDataProvider() *MoexDataProvider {
	m := MoexDataProvider{}
	m.baseUrl = "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities/"
	m.client = http.Client{
		Timeout: 3 * time.Second,
	}
	return &m
}

func (m *MoexDataProvider) GetStockPrice(ticker string, daysBackwards int, interval domain.Period) ([]domain.Candle, error) {
	if daysBackwards > 500 {
		return nil, fmt.Errorf("MOEX API doesn't support more than 500 days. Use pagination")
	}

	now := time.Now()
	from := now.AddDate(0, 0, -daysBackwards).Format("2006-01-02")
	till := now.Format("2006-01-02")
	url := fmt.Sprintf("%s%s/candles.json?from=%s&till=%s&interval=%d&iss.meta=off",
		m.baseUrl, ticker, from, till, int(interval))
	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response domain.CandlesApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	candles := domain.ParseCandles(response.Candles.Data)
	return candles, nil
}

func (m *MoexDataProvider) GetStockInfo(ticker string) (*domain.StockInfo, error) {
	url := fmt.Sprintf("https://iss.moex.com/iss/securities/%s.json?iss.meta=off&iss.only=description&description.columns=name,title,value", ticker)
	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MOEX API returned status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response domain.StockInfoApiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	stockInfo, err := domain.ParseStockInfo(response.Description.Data)
	if err != nil {
		return nil, err
	}

	return stockInfo, nil
}

func (m *MoexDataProvider) GetMarketCap(ticker string) (float64, error) {
	price, err := m.GetStockPrice(ticker, 5, domain.Period(60))
	if err != nil {
		return 0, err
	}

	stockInfo, err := m.GetStockInfo(ticker)
	if err != nil {
		return 0, err
	}

	return price[len(price) - 1].Close * float64(stockInfo.NumberOfShares), nil
}
