package domain

import "context"

type GeminiService interface {
	// Фундаментальный анализ компании на основе финансовой отчётности и рыночных данных
	AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period ReportPeriod) (string, error)

	// Функция, для работы с чатом
	GetChatResponse(ctx context.Context, prompt string, chatContext ChatContext) (string, error)

	// Функция для анализа сектора
	AnalyzeSector(ctx context.Context, sectorId int) (string, error)

	// Получить выжимку из старых отчетов для анализа нового
	GetCompanyHistory(ctx context.Context, ticker string) (string, error)

	ExtractResultFromReport(ctx context.Context, ticker string, year int, period ReportPeriod) (*ReportResults, error)
}
