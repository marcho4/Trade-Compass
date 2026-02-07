package application

import "context"

type Analyzer interface {
	GetReportS3Link(ctx context.Context, ticker string, period string) (string, error)
	GetContextToAnalyze(ctx context.Context, ticker string, period string) (string, error)
	GetAnalyzePrompt(ctx context.Context, ticker string, period string) (string, error)
	GenerateAnalysis(ctx context.Context, ticker string, period string) (string, error)
	SaveAnalysis(ctx context.Context, ticker string, period string, analysis string) error
}
