package parser

import (
	"ai-service/domain"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetReportS3Path(ctx context.Context, ticker, period string, year int) (string, error) {
	url := fmt.Sprintf("%s/reports/%s", c.baseURL, ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call parser API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("parser API returned status %d", resp.StatusCode)
	}

	var reports []domain.Report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return "", fmt.Errorf("failed to decode parser response: %w", err)
	}

	if year > 0 {
		for _, r := range reports {
			if r.Period == period && r.Year == year {
				return r.S3Path, nil
			}
		}
		return "", fmt.Errorf("report not found for %s year=%d period=%s", ticker, year, period)
	}

	var matching []domain.Report
	for _, r := range reports {
		if r.Period == period {
			matching = append(matching, r)
		}
	}

	if len(matching) == 0 {
		return "", fmt.Errorf("no reports found for %s period=%s", ticker, period)
	}

	sort.Slice(matching, func(i, j int) bool {
		return matching[i].Year > matching[j].Year
	})

	return matching[0].S3Path, nil
}

func (c *Client) GetLatestReportYear(ctx context.Context, ticker, period string) (int, error) {
	url := fmt.Sprintf("%s/reports/%s", c.baseURL, ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to call parser API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("parser API returned status %d", resp.StatusCode)
	}

	var reports []domain.Report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return 0, fmt.Errorf("failed to decode parser response: %w", err)
	}

	latestYear := 0
	for _, r := range reports {
		if r.Period == period && r.Year > latestYear {
			latestYear = r.Year
		}
	}

	if latestYear == 0 {
		return 0, fmt.Errorf("no reports found for %s period=%s", ticker, period)
	}

	return latestYear, nil
}
