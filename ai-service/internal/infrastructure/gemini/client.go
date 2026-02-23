package gemini

import (
	"ai-service/internal/domain"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/genai"
)

type Client struct {
	client *genai.Client
}

func NewClient(apiKey string, proxyURL string) (*Client, error) {
	cfg := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	if proxyURL != "" {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy URL: %w", err)
		}

		cfg.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parsedURL),
			},
		}
	}

	client, err := genai.NewClient(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) AnalyzeWithPDF(ctx context.Context, pdfBytes []byte, systemPrompt string) (string, error) {
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromBytes(pdfBytes, "application/pdf"),
			},
		},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				genai.NewPartFromText(systemPrompt),
			},
		},
	}

	result, err := c.client.Models.GenerateContent(ctx, string(domain.Pro), contents, config)
	if err != nil {
		return "", fmt.Errorf("gemini API call failed: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}

func (c *Client) GenerateText(ctx context.Context, prompt string, model domain.AIModel) (string, error) {
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{genai.NewPartFromText(prompt)},
		},
	}
	result, err := c.client.Models.GenerateContent(ctx, string(model), contents, nil)
	if err != nil {
		return "", fmt.Errorf("gemini API call failed: %w", err)
	}
	return strings.TrimSpace(result.Text()), nil
}
