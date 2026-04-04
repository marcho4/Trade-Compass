package parser

import (
	"ai-service/internal/domain/entity"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type reportsResponse struct {
	Reports []entity.Report `json:"reports"`
}

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

	var reports []entity.Report
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

	var matching []entity.Report
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

	var reports []entity.Report
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

func (c *Client) GetReports(ctx context.Context, ticker string) ([]entity.Report, error) {
	url := fmt.Sprintf("%s/reports/%s", c.baseURL, ticker)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call parser API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("parser API returned status %d", resp.StatusCode)
	}

	var body reportsResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode parser response: %w", err)
	}

	return body.Reports, nil
}

func (c *Client) IsLatestReport(ctx context.Context, ticker string, year int, periodMonths int) (bool, error) {
	reports, err := c.GetReports(ctx, ticker)
	if err != nil {
		return false, err
	}

	for _, r := range reports {
		rPeriod, _ := strconv.Atoi(r.Period)
		if r.Year > year || (r.Year == year && rPeriod > periodMonths) {
			return false, nil
		}
	}

	return true, nil
}
