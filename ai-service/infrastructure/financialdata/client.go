package financialdata

import (
	"ai-service/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetDraft(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	url := fmt.Sprintf("%s/raw-data/%s/draft?year=%d&period=%s", c.baseURL, ticker, year, period)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call financial-data API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("financial-data API returned status %d", resp.StatusCode)
	}

	var rawData domain.RawData
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &rawData, nil
}

func (c *Client) SaveDraft(ctx context.Context, rawData *domain.RawData) error {
	body, err := json.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %w", err)
	}

	url := fmt.Sprintf("%s/raw-data/%s", c.baseURL, rawData.Ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call financial-data API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("financial-data API returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) UpdateDraft(ctx context.Context, rawData *domain.RawData) error {
	body, err := json.Marshal(rawData)
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %w", err)
	}

	url := fmt.Sprintf("%s/raw-data/%s?year=%d&period=%s", c.baseURL, rawData.Ticker, rawData.Year, rawData.Period)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call financial-data API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("financial-data API returned status %d", resp.StatusCode)
	}

	return nil
}
