package infrastructure

import (
	"context"
	"encoding/json"
	"financial_data/internal/domain"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	PRICE_TTL   = 15 * time.Minute
	COMPANY_TTL = 24 * time.Hour
)

type MoexDataProvider struct {
	baseUrl string
	client  http.Client
	redis   *redis.Client
}

func NewMoexDataProvider(redis *redis.Client) *MoexDataProvider {
	m := MoexDataProvider{}
	m.baseUrl = "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities/"
	m.client = http.Client{
		Timeout: 3 * time.Second,
	}
	m.redis = redis
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

	cache, err := m.redis.Get(context.TODO(), url).Result()
	if err != nil && err != redis.Nil {
		slog.Warn("redis get error", slog.String("key", url), slog.Any("err", err))
	}
	if err == nil {
		if candles, err := unmarshalCandles([]byte(cache)); err == nil {
			return candles, nil
		}
		slog.Warn("failed to unmarshal cached candles, fetching from API", slog.String("ticker", ticker))
	}

	resp, err := m.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	candles, err := unmarshalCandles(body)
	if err != nil {
		return nil, err
	}

	if err := m.redis.Set(context.TODO(), url, string(body), PRICE_TTL).Err(); err != nil {
		slog.Warn("failed to cache stock price", slog.String("ticker", ticker), slog.Any("err", err))
	}

	return candles, nil
}

func (m *MoexDataProvider) GetStockInfo(ticker string) (*domain.StockInfo, error) {
	url := fmt.Sprintf("https://iss.moex.com/iss/securities/%s.json?iss.meta=off&iss.only=description&description.columns=name,title,value", ticker)

	cache, err := m.redis.Get(context.TODO(), url).Result()
	if err != nil && err != redis.Nil {
		slog.Warn("redis get error", slog.String("key", url), slog.Any("err", err))
	}
	if err == nil {
		if stockInfo, err := unmarshalStockInfo([]byte(cache)); err == nil {
			return stockInfo, nil
		}
		slog.Warn("failed to unmarshal cached stock info, fetching from API", slog.String("ticker", ticker))
	}

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

	stockInfo, err := unmarshalStockInfo(body)
	if err != nil {
		return nil, err
	}

	if err := m.redis.Set(context.TODO(), url, string(body), COMPANY_TTL).Err(); err != nil {
		slog.Warn("failed to cache stock info", slog.String("ticker", ticker), slog.Any("err", err))
	}

	return stockInfo, nil
}

func unmarshalCandles(data []byte) ([]domain.Candle, error) {
	var response domain.CandlesApiResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return domain.ParseCandles(response.Candles.Data), nil
}

func unmarshalStockInfo(data []byte) (*domain.StockInfo, error) {
	var response domain.StockInfoApiResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	return domain.ParseStockInfo(response.Description.Data)
}

func (m *MoexDataProvider) GetPriceAt(ticker string, date time.Time) (float64, error) {
	from := date.AddDate(0, 0, -7).Format("2006-01-02")
	till := date.Format("2006-01-02")
	url := fmt.Sprintf("%s%s/candles.json?from=%s&till=%s&interval=24&iss.meta=off",
		m.baseUrl, ticker, from, till)

	cache, err := m.redis.Get(context.TODO(), url).Result()
	if err != nil && err != redis.Nil {
		slog.Warn("redis get error", slog.String("key", url), slog.Any("err", err))
	}
	if err == nil {
		if candles, err := unmarshalCandles([]byte(cache)); err == nil && len(candles) > 0 {
			return candles[len(candles)-1].Close, nil
		}
	}

	resp, err := m.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	candles, err := unmarshalCandles(body)
	if err != nil {
		return 0, err
	}
	if len(candles) == 0 {
		return 0, fmt.Errorf("no candles found for %s at %s", ticker, till)
	}

	if err := m.redis.Set(context.TODO(), url, string(body), PRICE_TTL).Err(); err != nil {
		slog.Warn("failed to cache price at date", slog.String("ticker", ticker), slog.Any("err", err))
	}

	return candles[len(candles)-1].Close, nil
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

	return price[len(price)-1].Close * float64(stockInfo.NumberOfShares), nil
}
