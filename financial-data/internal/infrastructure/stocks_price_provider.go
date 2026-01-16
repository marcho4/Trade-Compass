package infrastructure

import (
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MoexPriceProvider struct {
	baseUrl string
	client  http.Client
}

func NewMoexPriceProvider() *MoexPriceProvider {
	m := MoexPriceProvider{}
	m.baseUrl = "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities/"
	m.client = http.Client{
		Timeout: 3 * time.Second,
	}
	return &m
}

func (m *MoexPriceProvider) GetStockPrice(ticker string, daysBackwards int, interval domain.Period) ([]domain.Candle, error) {
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

func (m *MoexPriceProvider) GetStockInfo(ticker string) {
	// url := fmt.Sprintf("%s%s/securities.json", m.baseUrl, ticker)
	// resp, err := m.client.Get(url)
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }
	// var response domain.StockInfoApiResponse
	// err = json.Unmarshal(body, &response)
	// if err != nil {
	// 	return nil, err
	// }
	// return &response.StockInfo, nil
}
