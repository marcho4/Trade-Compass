package usecase

import (
	"ai-service/internal/domain/entity"
	"context"
	"fmt"
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
	ai                AIService
	finData           FinancialDataGateway
	riskAndGrowthRepo RiskAndGrowthRepository
}

func NewScenarioGenerator(
	ai AIService,
	finData FinancialDataGateway,
	riskAndGrowthRepo RiskAndGrowthRepository,
) *ScenarioGenerator {
	return &ScenarioGenerator{
		ai:                ai,
		finData:           finData,
		riskAndGrowthRepo: riskAndGrowthRepo,
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

	latest, ok := getLatestFullYearData(history)
	if !ok {
		return fmt.Errorf("no confirmed annual data found for %s", task.Ticker)
	}

	wacc := calculateWACC(latest, cbRate)

	scenarios, err := s.ai.GenerateScenarios(ctx, task.Ticker, YearsToForecast, history, cbRate, wacc, riskAndGrowth)
	if err != nil {
		return fmt.Errorf("generate scenarios: %w", err)
	}

	dcfInput := buildDCFInput(latest, wacc)

	dcfResult := Calculate(dcfInput, scenarios)

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
func calculateWACC(d entity.RawData, cbRate *entity.CBRate) float64 {
	rf := cbRate.Rate / 100
	ke := rf + defaultBeta*equityRiskPremium

	var kd float64
	if d.Debt != nil && *d.Debt != 0 && d.InterestOnLoans != nil && *d.InterestOnLoans != 0 {
		kd = math.Abs(float64(*d.InterestOnLoans)) / float64(*d.Debt)
	} else {
		kd = rf + 0.02 // fallback: rf + кредитный спред
	}

	taxRate := effectiveTaxRate(d)

	var e, debt float64
	if d.MarketCap != nil {
		e = float64(*d.MarketCap)
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

func buildDCFInput(d entity.RawData, wacc float64) entity.DCFInput {
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
	if d.SharesOutstanding != nil {
		input.SharesOutstanding = *d.SharesOutstanding
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
