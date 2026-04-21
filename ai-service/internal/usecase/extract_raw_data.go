package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"time"

	docs "ai-service/internal/docs"
	"ai-service/internal/domain/entity"
)

type ExtractRawDataUsecase struct {
	ai        AIProvider
	fd        FinancialDataGateway
	parser    ParserGateway
	publisher MessagePublisher
	storage   StorageClient
}

func NewExtractRawDataUsecase(
	ai AIProvider,
	fd FinancialDataGateway,
	parser ParserGateway,
	publisher MessagePublisher,
	storage StorageClient,
) *ExtractRawDataUsecase {
	return &ExtractRawDataUsecase{
		ai:        ai,
		fd:        fd,
		parser:    parser,
		publisher: publisher,
		storage:   storage,
	}
}

func (u *ExtractRawDataUsecase) Execute(ctx context.Context, task entity.Task) error {
	logger := slog.With(
		slog.String("id", task.Id),
		slog.String("ticker", task.Ticker),
		slog.Int("year", task.Year),
		slog.String("period", task.Period),
	)

	period := entity.ReportPeriod(task.Period)

	logger.Info("extracting raw data")

	existing, err := u.fd.GetRawData(ctx, task.Ticker, task.Year, period)
	if err != nil {
		return fmt.Errorf("check existing raw data: %w", err)
	}

	if existing == nil {
		prompt := docs.RawDataAgentPrompt() + "\n## Ticker\n\n" + task.Ticker

		pdfBytes, err := u.storage.DownloadPDF(ctx, task.ReportURL)
		if err != nil {
			return fmt.Errorf("download PDF: %w", err)
		}

		var temp float32 = 0.1

		text, err := u.ai.AnalyzeWithPDF(
			ctx,
			pdfBytes,
			prompt,
			entity.Pro,
			GenerateParams{
				Temperature:    &temp,
				ResponseSchema: getResponseSchema(),
			},
		)

		if err != nil {
			return fmt.Errorf("ai call with pdf: %w", err)
		}

		prompt = docs.RawDataValidatorPrompt()
		prompt += "\n\n## Received raw data from an agent\n\n" + text

		validationResults, err := u.ai.AnalyzeWithPDF(
			ctx,
			pdfBytes,
			prompt,
			entity.Pro,
			GenerateParams{
				Temperature:    &temp,
				ResponseSchema: getValidatorResponseSchema(),
			},
		)

		if err != nil {
			return fmt.Errorf("validate: %w", err)
		}

		logger.Info("Validation Result", "res", validationResults)

		prompt = docs.RawDataAgentPrompt() +
			"\n## Ticker\n\n" + task.Ticker +
			"\n\n## Validation result\n\nFix every issue with errorLevel \"critical\" or \"high\". Issues with errorLevel \"warning\" should be reviewed but fixed only if clearly wrong.\n\n" +
			validationResults +
			"\n\n## Old result\n\n" + text

		finalRawData, err := u.ai.AnalyzeWithPDF(
			ctx,
			pdfBytes,
			prompt,
			entity.Pro,
			GenerateParams{
				Temperature:    &temp,
				ResponseSchema: getResponseSchema(),
			},
		)

		if err != nil {
			return fmt.Errorf("ai call with pdf: %w", err)
		}

		var rawData entity.RawData
		if err := json.Unmarshal([]byte(finalRawData), &rawData); err != nil {
			return fmt.Errorf("unmarshal raw data: %w", err)
		}

		periodEnd := periodEndDate(task.Year, period)
		stockInfo, stockErr := u.fd.GetStockInfo(ctx, task.Ticker)
		price, priceErr := u.fd.GetPriceAt(ctx, task.Ticker, periodEnd)

		if stockErr != nil || priceErr != nil {
			logger.Warn(
				"failed to get stock info or price, skipping market cap calculation",
				slog.Any("stock error", stockErr),
				slog.Any("price err", priceErr),
			)
		} else {
			unitDivisor := unitDivisorForReportUnits(rawData.ReportUnits)
			shares := int64(stockInfo.NumberOfShares)
			rawData.SharesOutstanding = &shares
			marketCap := int64(math.Round(price*float64(stockInfo.NumberOfShares))) / unitDivisor
			rawData.MarketCap = &marketCap
		}

		rawData.ComputeDerivedFields()
		rawData.Status = entity.RawDataStatusConfirmed

		if err := u.fd.SaveDraft(ctx, &rawData); err != nil {
			return fmt.Errorf("save raw data: %w", err)
		}
	}

	nextTask := entity.Task{
		Id:        task.Id,
		Ticker:    task.Ticker,
		Year:      task.Year,
		Period:    task.Period,
		ReportURL: task.ReportURL,
		Type:      entity.RawDataSuccess,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal analyze task: %w", err)
	}

	if err := u.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish analyze task: %w", err)
	}

	logger.Info("raw data extraction succeed")

	return nil
}

func unitDivisorForReportUnits(units string) int64 {
	switch units {
	case "thousands":
		return 1_000
	case "millions":
		return 1_000_000
	case "billions":
		return 1_000_000_000
	default:
		return 1
	}
}

func periodEndDate(year int, period entity.ReportPeriod) time.Time {
	months, ok := entity.PeriodToMonths[string(period)]
	if !ok {
		months = 12
	}
	lastDay := time.Date(year, time.Month(months)+1, 0, 0, 0, 0, 0, time.UTC)
	return lastDay
}

func getResponseSchema() *Schema {
	str := &Schema{Type: TypeString}
	integer := &Schema{Type: TypeInteger}
	warnings := &Schema{Type: TypeArray, Items: str}

	properties := map[string]*Schema{
		"ticker":      str,
		"year":        integer,
		"period":      {Type: TypeString, Enum: []string{"Q1", "Q2", "Q3", "YEAR"}},
		"status":      str,
		"reportUnits": {Type: TypeString, Enum: []string{"units", "thousands", "millions", "billions"}},
		"companyType": {Type: TypeString, Enum: []string{"bank", "corporate"}},

		// Income Statement
		"revenue":           integer,
		"costOfRevenue":     integer,
		"grossProfit":       integer,
		"operatingExpenses": integer,
		"otherIncome":       integer,
		"otherExpenses":     integer,
		"ebit":              integer,
		"interestIncome":    integer,
		"interestExpense":   integer,
		"profitBeforeTax":   integer,
		"taxExpense":        integer,
		"netProfit":         integer,
		"netProfitParent":   integer,

		// Balance Sheet - Assets
		"totalAssets":           integer,
		"currentAssets":         integer,
		"cashAndEquivalents":    integer,
		"inventories":           integer,
		"receivables":           integer,
		"fixedAssets":           integer,
		"rightOfUseAssets":      integer,
		"intangibleAssets":      integer,
		"goodwill":              integer,
		"totalNonCurrentAssets": integer,

		// Balance Sheet - Liabilities & Equity
		"totalLiabilities":   integer,
		"currentLiabilities": integer,
		"longTermDebt":       integer,
		"shortTermDebt":      integer,
		"ltLeaseLiabilities": integer,
		"stLeaseLiabilities": integer,
		"tradePayables":      integer,
		"equity":             integer,
		"equityParent":       integer,
		"treasuryShares":     integer,
		"retainedEarnings":   integer,

		// Cash Flow
		"operatingCashFlow": integer,
		"investingCashFlow": integer,
		"financingCashFlow": integer,
		"daFixedRou":        integer,
		"daIntangibles":     integer,
		"capexFa":           integer,
		"capexIa":           integer,
		"dividendsPaid":     integer,
		"leasePayments":     integer,
		"acquisitionsNet":   integer,
		"interestPaid":      integer,
		"debtProceeds":      integer,
		"debtRepayments":    integer,

		// Interest breakdown
		"interestOnLeases": integer,
		"interestOnLoans":  integer,

		// Bank-specific
		"netInterestIncome":    integer,
		"commissionIncome":     integer,
		"commissionExpense":    integer,
		"netCommissionIncome":  integer,
		"creditLossProvision":  integer,
		"interbankLiabilities": integer,

		"warnings": warnings,
	}

	return &Schema{
		Type:       TypeObject,
		Properties: properties,
		Required: []string{
			"ticker", "year", "period", "reportUnits", "companyType",
			"revenue", "costOfRevenue", "grossProfit", "operatingExpenses",
			"ebit", "interestExpense", "profitBeforeTax", "taxExpense", "netProfit",
			"totalAssets", "currentAssets", "cashAndEquivalents", "inventories", "receivables",
			"totalLiabilities", "currentLiabilities", "longTermDebt", "shortTermDebt",
			"equity", "retainedEarnings",
			"operatingCashFlow", "investingCashFlow", "financingCashFlow",
			"warnings",
		},
	}
}

func getValidatorResponseSchema() *Schema {
	str := &Schema{Type: TypeString}

	properties := map[string]*Schema{
		"rule":       str,
		"errorLevel": {Type: TypeString, Enum: []string{"critical", "high", "warning"}},
		"fieldName":  str,
		"reason":     str,
		"hint":       str,
	}

	errReport := Schema{
		Type:       TypeObject,
		Properties: properties,
		Required:   []string{"fieldName", "errorLevel", "reason", "rule", "hint"},
	}

	validatorReport := Schema{
		Type:  TypeArray,
		Items: &errReport,
	}

	return &validatorReport
}
