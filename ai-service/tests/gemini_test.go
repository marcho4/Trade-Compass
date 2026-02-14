package tests

import (
	"ai-service/internal/infrastructure/gemini"
	"context"
	"os"
	"strings"
	"testing"
)

func GeminiTest(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY is not set")
	}

	ctx := context.Background()
	client, err := gemini.NewClient(apiKey)
	if err != nil {
		t.Fatalf("failed to create gemini client: %v", err)
	}

	result, err := client.GenerateText(ctx, "Explain how AI works in a few words")
	if err != nil {
		t.Fatalf("failed to call gemini API: %v", err)
	}

	if strings.TrimSpace(result) == "" {
		t.Fatal("gemini returned empty response")
	}

	t.Logf("Gemini response: %s", result)
}
