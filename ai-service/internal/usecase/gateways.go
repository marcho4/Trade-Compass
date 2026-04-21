package usecase

import (
	"context"
	"time"

	"ai-service/internal/domain/entity"
)

type FinancialDataGateway interface {
	GetDailyPrices(ctx context.Context, ticker string) ([]entity.Candle, error)
	GetCBRates(ctx context.Context) (*entity.CBRate, error)
	GetMarketCap(ctx context.Context, ticker string) (float64, error)
	GetPriceAt(ctx context.Context, ticker string, date time.Time) (float64, error)
	GetStockInfo(ctx context.Context, ticker string) (*entity.StockInfo, error)
	GetRawData(ctx context.Context, ticker string, year int, period entity.ReportPeriod) (*entity.RawData, error)
	GetRawDataHistory(ctx context.Context, ticker string) ([]entity.RawData, error)
	SaveDraft(ctx context.Context, rawData *entity.RawData) error
}

type ParserGateway interface {
	IsLatestReport(ctx context.Context, ticker string, year, periodMonths int) (bool, error)
	GetLatestReport(ctx context.Context, ticker string) (*entity.Report, error)
}

type MessagePublisher interface {
	PublishMessage(ctx context.Context, value []byte) error
}

type StorageClient interface {
	DownloadPDF(ctx context.Context, url string) ([]byte, error)
}

type AIProvider interface {
	AnalyzeWithPDF(ctx context.Context, pdfBytes []byte, systemPrompt string, model entity.AIModel, params GenerateParams) (string, error)
	GenerateText(ctx context.Context, prompt string, model entity.AIModel, params GenerateParams) (string, error)
}

type GenerateParams struct {
	Temperature    *float32
	GoogleSearch   bool
	ResponseSchema *Schema
}

type SchemaType string

const (
	TypeObject  SchemaType = "object"
	TypeArray   SchemaType = "array"
	TypeString  SchemaType = "string"
	TypeNumber  SchemaType = "number"
	TypeInteger SchemaType = "integer"
	TypeBoolean SchemaType = "boolean"
)

type Schema struct {
	Type       SchemaType
	Properties map[string]*Schema
	Items      *Schema
	Enum       []string
	Required   []string
}

func Float32Ptr(v float32) *float32 { return &v }
