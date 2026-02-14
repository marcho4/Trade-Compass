package domain

import "context"

type GeminiService interface {
	// Фундаментальный анализ компании на основе финансовой отчётности и рыночных данных
	AnalyzeReport(ctx context.Context, ticker, reportUrl string, year int, period ReportPeriod) (string, error)

	// Функция, для работы с чатом
	GetChatResponse(prompt string, chatContext ChatContext) (string, error)

	// Функция для анализа сектора
	AnalyzeSector(sectorId int) (string, error)

	// Получить выжимку из старых отчетов для анализа нового
	GetCompanyHistory(ticker string) (string, error)

	// Получить сырые данные из отчета
	ExtractDataFromReport(reportUrl string) (string, error)
}
