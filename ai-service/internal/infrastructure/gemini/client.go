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
			return nil, fmt.Errorf("parse proxy URL: %w", err)
		}

		cfg.HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(parsedURL),
			},
		}
	}

	client, err := genai.NewClient(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) AnalyzeWithPDF(ctx context.Context, pdfBytes []byte, systemPrompt string, model domain.AIModel) (string, error) {
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

	result, err := c.client.Models.GenerateContent(ctx, string(model), contents, config)
	if err != nil {
		return "", fmt.Errorf("gemini API call failed: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}

type GenerateOption func(*genai.GenerateContentConfig)

func WithGoogleSearch() GenerateOption {
	return func(cfg *genai.GenerateContentConfig) {
		cfg.Tools = append(cfg.Tools, &genai.Tool{
			GoogleSearch: &genai.GoogleSearch{},
		})
	}
}

func WithResponseSchema(schema *genai.Schema) GenerateOption {
	return func(cfg *genai.GenerateContentConfig) {
		cfg.ResponseMIMEType = "application/json"
		cfg.ResponseSchema = schema
	}
}

func (c *Client) GenerateText(ctx context.Context, prompt string, model domain.AIModel, opts ...GenerateOption) (string, error) {
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{genai.NewPartFromText(prompt)},
		},
	}

	var config *genai.GenerateContentConfig
	if len(opts) > 0 {
		config = &genai.GenerateContentConfig{}
		for _, opt := range opts {
			opt(config)
		}
	}

	result, err := c.client.Models.GenerateContent(ctx, string(model), contents, config)
	if err != nil {
		return "", fmt.Errorf("gemini API call failed: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}
