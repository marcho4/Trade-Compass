package gemini

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"ai-service/internal/domain/entity"
	"ai-service/internal/usecase"

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

	return &Client{client: client}, nil
}

func (c *Client) AnalyzeWithPDF(ctx context.Context, pdfBytes []byte, prompt string, model entity.AIModel, params usecase.GenerateParams) (string, error) {
	contents := []*genai.Content{
		{
			Role: "user",
			Parts: []*genai.Part{
				genai.NewPartFromBytes(pdfBytes, "application/pdf"),
			},
		},
		{
			Role:  "user",
			Parts: []*genai.Part{genai.NewPartFromText(prompt)},
		},
	}

	config := buildConfig(params)

	var thinkingBudget int32 = 16000

	thinkingConfig := &genai.ThinkingConfig{}
	thinkingConfig.ThinkingBudget = &thinkingBudget
	thinkingConfig.IncludeThoughts = false

	config.ThinkingConfig = thinkingConfig

	result, err := c.client.Models.GenerateContent(ctx, string(model), contents, config)
	if err != nil {
		return "", fmt.Errorf("call gemini: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}

func (c *Client) GenerateText(ctx context.Context, prompt string, model entity.AIModel, params usecase.GenerateParams) (string, error) {
	contents := []*genai.Content{
		{
			Role:  "user",
			Parts: []*genai.Part{genai.NewPartFromText(prompt)},
		},
	}

	config := buildConfig(params)

	result, err := c.client.Models.GenerateContent(ctx, string(model), contents, config)
	if err != nil {
		return "", fmt.Errorf("call gemini: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}

func buildConfig(params usecase.GenerateParams) *genai.GenerateContentConfig {
	cfg := &genai.GenerateContentConfig{}

	if params.Temperature != nil {
		cfg.Temperature = genai.Ptr(*params.Temperature)
	}

	if params.GoogleSearch {
		cfg.Tools = append(cfg.Tools, &genai.Tool{
			GoogleSearch: &genai.GoogleSearch{},
		})
	}

	if params.ResponseSchema != nil {
		cfg.ResponseMIMEType = "application/json"
		cfg.ResponseSchema = convertSchema(params.ResponseSchema)
	}

	return cfg
}

func convertSchema(s *usecase.Schema) *genai.Schema {
	if s == nil {
		return nil
	}

	out := &genai.Schema{
		Type:     toGenaiType(s.Type),
		Enum:     s.Enum,
		Required: s.Required,
	}

	if len(s.Properties) > 0 {
		out.Properties = make(map[string]*genai.Schema, len(s.Properties))
		for k, v := range s.Properties {
			out.Properties[k] = convertSchema(v)
		}
	}

	if s.Items != nil {
		out.Items = convertSchema(s.Items)
	}

	return out
}

func toGenaiType(t usecase.SchemaType) genai.Type {
	switch t {
	case usecase.TypeObject:
		return genai.TypeObject
	case usecase.TypeArray:
		return genai.TypeArray
	case usecase.TypeString:
		return genai.TypeString
	case usecase.TypeNumber:
		return genai.TypeNumber
	case usecase.TypeInteger:
		return genai.TypeInteger
	case usecase.TypeBoolean:
		return genai.TypeBoolean
	default:
		return genai.TypeUnspecified
	}
}
