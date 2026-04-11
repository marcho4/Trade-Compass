package usecase

import (
	docs "ai-service/internal/docs"
	"ai-service/internal/domain/entity"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"
)

const (
	YearsToForecast   = 3
	equityRiskPremium = 0.07 // ERP для российского рынка
	defaultBeta       = 1.0
	defaultTaxRate    = 0.25
)

type ScenarioGenerator struct {
	ai                AIProvider
	finData           FinancialDataGateway
	parserGateway     ParserGateway
	riskAndGrowthRepo RiskAndGrowthRepository
	scenarioRepo      ScenarioRepository
	dcfRepo           DCFResultsRepository
	transactor        Transactor
	publisher         MessagePublisher
}

func NewScenarioGenerator(
	ai AIProvider,
	finData FinancialDataGateway,
	parserGateway ParserGateway,
	riskAndGrowthRepo RiskAndGrowthRepository,
	scenarioRepo ScenarioRepository,
	dcfRepo DCFResultsRepository,
	transactor Transactor,
	publisher MessagePublisher,
) *ScenarioGenerator {
	return &ScenarioGenerator{
		ai:                ai,
		finData:           finData,
		parserGateway:     parserGateway,
		riskAndGrowthRepo: riskAndGrowthRepo,
		scenarioRepo:      scenarioRepo,
		dcfRepo:           dcfRepo,
		transactor:        transactor,
		publisher:         publisher,
	}
}

func (s *ScenarioGenerator) Execute(ctx context.Context, task entity.Task) error {
	history, err := s.finData.GetRawDataHistory(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get history: %w", err)
	}

	cbRate, err := s.finData.GetCBRates(ctx)
	if err != nil {
		return fmt.Errorf("get cb rate: %w", err)
	}

	riskAndGrowth, err := s.riskAndGrowthRepo.GetFreshRiskAndGrowth(ctx, task.Ticker, 72*time.Hour)
	if err != nil {
		return fmt.Errorf("get risk and growth: %w", err)
	}

	if riskAndGrowth == nil {
		return errors.New("riskAndGrowth is nil")
	}

	latest, ok := getLatestFullYearData(history)
	if !ok {
		return fmt.Errorf("no confirmed annual data found for %s", task.Ticker)
	}

	stockInfo, err := s.finData.GetStockInfo(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get stock info: %w", err)
	}

	marketCap, err := s.finData.GetMarketCap(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get market cap: %w", err)
	}

	wacc := calculateWACC(latest, cbRate, marketCap)

	historyJSON, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}

	prompt := docs.ScenarioGeneratorPrompt()

	prompt += fmt.Sprintf("\n\n## Кол-во лет\n\n%d", YearsToForecast)

	prompt += fmt.Sprintf("\n\n## Исторические данные компании\n\nТикер: %s\n\n%s", task.Ticker, string(historyJSON))

	prompt += fmt.Sprintf("\n\n## Макроэкономические данные\n\nСтавка ЦБ РФ: %.2f%%\nWACC: %.4f", cbRate.Rate, wacc)

	var risks, growthFactors []entity.RiskAndGrowthFactor
	for _, f := range riskAndGrowth.Factors {
		if f.Type == entity.FactorRisk {
			risks = append(risks, f)
		} else {
			growthFactors = append(growthFactors, f)
		}
	}

	risksJSON, _ := json.Marshal(risks)
	growthJSON, _ := json.Marshal(growthFactors)

	prompt += fmt.Sprintf("\n\n## Факторы риска\n\n%s", string(risksJSON))

	prompt += fmt.Sprintf("\n\n## Факторы роста\n\n%s", string(growthJSON))

	factorSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"factor": {Type: TypeString},
			"impact": {Type: TypeString},
		},
		Required: []string{"factor", "impact"},
	}

	assumptionSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"year":              {Type: TypeInteger},
			"revenue_growth":    {Type: TypeNumber},
			"cogs_pct_revenue":  {Type: TypeNumber},
			"sga_pct_revenue":   {Type: TypeNumber},
			"tax_rate":          {Type: TypeNumber},
			"capex_pct_revenue": {Type: TypeNumber},
			"da_pct_revenue":    {Type: TypeNumber},
			"nwc_pct_revenue":   {Type: TypeNumber},
		},
		Required: []string{"year", "revenue_growth", "cogs_pct_revenue", "sga_pct_revenue", "tax_rate", "capex_pct_revenue", "da_pct_revenue", "nwc_pct_revenue"},
	}

	scenarioSchema := &Schema{
		Type: TypeObject,
		Properties: map[string]*Schema{
			"id":                     {Type: TypeString},
			"name":                   {Type: TypeString},
			"description":            {Type: TypeString},
			"probability":            {Type: TypeNumber},
			"terminal_growth_rate":   {Type: TypeNumber},
			"growth_factors_applied": {Type: TypeArray, Items: factorSchema},
			"risks_applied":          {Type: TypeArray, Items: factorSchema},
			"assumptions":            {Type: TypeArray, Items: assumptionSchema},
		},
		Required: []string{"id", "name", "description", "probability", "terminal_growth_rate", "assumptions"},
	}

	text, err := s.ai.GenerateText(ctx, prompt, entity.Pro, GenerateParams{
		ResponseSchema: &Schema{
			Type:  TypeArray,
			Items: scenarioSchema,
		},
	})
	if err != nil {
		return fmt.Errorf("generate scenarios: %w", err)
	}

	var dtos []scenarioDTO
	if err := json.Unmarshal([]byte(text), &dtos); err != nil {
		slog.Error("failed to parse scenarios response", slog.String("ai_response", text))
		return fmt.Errorf("parse scenarios response: %w", err)
	}
	scenarios := mapScenariosToDomain(dtos)

	dcfInput := buildDCFInput(latest, wacc, stockInfo.NumberOfShares)

	dcfResult := Calculate(dcfInput, scenarios)
	dcfResult.ID = task.Id

	if err := s.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		if err := s.scenarioRepo.SaveScenarios(txCtx, task.Id, task.Ticker, scenarios); err != nil {
			return fmt.Errorf("save scenarios: %w", err)
		}
		if err := s.dcfRepo.SaveDCFResults(txCtx, task.Ticker, dcfResult); err != nil {
			return fmt.Errorf("save dcf results: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("save results: %w", err)
	}

	latestReport, err := s.parserGateway.GetLatestReport(ctx, task.Ticker)
	if err != nil {
		return fmt.Errorf("get latest report: %w", err)
	}

	reportPeriod, ok := entity.MonthsToPeriod[latestReport.Period]
	if !ok {
		return fmt.Errorf("unknown period months %q for ticker %s", latestReport.Period, task.Ticker)
	}

	nextTask := entity.Task{
		Id:        task.Id,
		Ticker:    task.Ticker,
		Year:      latestReport.Year,
		Period:    string(reportPeriod),
		ReportURL: latestReport.S3Path,
		Type:      entity.Analyze,
	}

	payload, err := json.Marshal(nextTask)
	if err != nil {
		return fmt.Errorf("marshal analyze task: %w", err)
	}

	if err := s.publisher.PublishMessage(ctx, payload); err != nil {
		return fmt.Errorf("publish analyze task: %w", err)
	}

	return nil
}

// getLatestFullYearData возвращает последние подтверждённые годовые данные из истории.
func getLatestFullYearData(history []entity.RawData) (entity.RawData, bool) {
	annual := make([]entity.RawData, 0, len(history))
	for _, d := range history {
		if d.Period == entity.YEAR && d.Status == entity.RawDataStatusConfirmed {
			annual = append(annual, d)
		}
	}
	if len(annual) == 0 {
		return entity.RawData{}, false
	}

	sort.Slice(annual, func(i, j int) bool {
		return annual[i].Year > annual[j].Year
	})
	return annual[0], true
}

// calculateWACC вычисляет WACC через CAPM для стоимости капитала и фактическую
// стоимость долга из отчётности. CB rate передаётся в долях (0.16 = 16%).
// marketCap — актуальная рыночная капитализация в рублях из MOEX API.
func calculateWACC(d entity.RawData, cbRate *entity.CBRate, marketCap float64) float64 {
	rf := cbRate.Rate / 100
	ke := rf + defaultBeta*equityRiskPremium

	var kd float64
	if d.Debt != nil && *d.Debt != 0 && d.InterestOnLoans != nil && *d.InterestOnLoans != 0 {
		kd = math.Abs(float64(*d.InterestOnLoans)) / float64(*d.Debt)
	} else {
		kd = rf + 0.02 // fallback: rf + кредитный спред
	}

	taxRate := effectiveTaxRate(d)

	divisor := unitDivisor(d.ReportUnits)
	var e, debt float64
	if marketCap > 0 {
		e = marketCap / divisor
	} else if d.Equity != nil {
		e = float64(*d.Equity)
	}
	if d.Debt != nil {
		debt = float64(*d.Debt)
	}

	total := e + debt
	if total == 0 {
		return ke
	}

	return ke*(e/total) + kd*(1-taxRate)*(debt/total)
}

func unitDivisor(reportUnits string) float64 {
	switch reportUnits {
	case "billions":
		return 1_000_000_000
	case "millions":
		return 1_000_000
	case "thousands":
		return 1_000
	default:
		return 1
	}
}

func buildDCFInput(d entity.RawData, wacc float64, numberOfShares int) entity.DCFInput {
	input := entity.DCFInput{WACC: wacc}

	if d.Revenue != nil {
		input.BaseRevenue = float64(*d.Revenue)
	}
	if d.WorkingCapital != nil {
		input.BaseNWC = float64(*d.WorkingCapital)
	}
	if d.NetDebt != nil {
		input.NetDebt = float64(*d.NetDebt)
	}

	divisor := unitDivisor(d.ReportUnits)
	if numberOfShares > 0 {
		input.SharesOutstanding = float64(numberOfShares) / divisor
	}

	return input
}

func effectiveTaxRate(d entity.RawData) float64 {
	if d.TaxExpense == nil || d.ProfitBeforeTax == nil || *d.ProfitBeforeTax == 0 {
		return defaultTaxRate
	}
	rate := math.Abs(float64(*d.TaxExpense)) / math.Abs(float64(*d.ProfitBeforeTax))
	if rate <= 0 || rate >= 1 {
		return defaultTaxRate
	}
	return rate
}
