package port

import (
	"ai-service/internal/domain/entity"
	"context"
)

type AIClient interface {
	AnalyzeWithPDF(ctx context.Context, pdfBytes []byte, systemPrompt string, model entity.AIModel) (string, error)
	GenerateText(ctx context.Context, prompt string, model entity.AIModel, opts ...GenerateOption) (string, error)
}

type GenerateOption func(cfg any)
