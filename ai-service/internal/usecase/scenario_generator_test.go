package usecase_test

import (
	"context"
	"testing"
	"time"

	"ai-service/internal/domain/entity"
	"ai-service/internal/usecase"
	"ai-service/internal/usecase/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptr64(v int64) *int64 { return &v }

func moexRawData() entity.RawData {
	return entity.RawData{
		Ticker:      "MOEX",
		Year:        2024,
		Period:      entity.YEAR,
		Status:      entity.RawDataStatusConfirmed,
		ReportUnits: "millions",

		Revenue:         ptr64(129_000),
		WorkingCapital:  ptr64(50_000),
		NetDebt:         ptr64(-670_000),
		Debt:            ptr64(5_000),
		InterestOnLoans: ptr64(-500),
		Equity:          ptr64(340_000),
		ProfitBeforeTax: ptr64(80_000),
		TaxExpense:      ptr64(-16_000),
	}
}

func TestBuildDCFInput_UnitsConversion(t *testing.T) {
	d := moexRawData()
	numberOfShares := 2_276_401_458
	wacc := 0.15

	input := usecase.BuildDCFInput(d, wacc, numberOfShares)

	assert.Equal(t, 129_000.0, input.BaseRevenue)
	assert.Equal(t, 50_000.0, input.BaseNWC)
	assert.Equal(t, -670_000.0, input.NetDebt)
	assert.Equal(t, wacc, input.WACC)

	expectedShares := float64(numberOfShares) / 1_000_000
	assert.InDelta(t, expectedShares, input.SharesOutstanding, 0.001)
}

func TestBuildDCFInput_ZeroShares(t *testing.T) {
	d := moexRawData()
	input := usecase.BuildDCFInput(d, 0.15, 0)
	assert.Equal(t, 0.0, input.SharesOutstanding)
}

func TestUnitDivisor(t *testing.T) {
	cases := []struct {
		units    string
		expected float64
	}{
		{"millions", 1_000_000},
		{"billions", 1_000_000_000},
		{"thousands", 1_000},
		{"units", 1},
		{"", 1},
	}
	for _, tc := range cases {
		assert.Equal(t, tc.expected, usecase.UnitDivisor(tc.units), "units=%q", tc.units)
	}
}

const mockScenarioJSON = `[
  {
    "id": "base",
    "name": "Базовый",
    "description": "Постепенное снижение ставки",
    "probability": 0.5,
    "terminal_growth_rate": 0.03,
    "assumptions": [
      {"year": 2025, "revenue_growth": 0.08, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.40, "tax_rate": 0.25, "capex_pct_revenue": 0.08, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.05},
      {"year": 2026, "revenue_growth": 0.06, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.40, "tax_rate": 0.25, "capex_pct_revenue": 0.08, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.05},
      {"year": 2027, "revenue_growth": 0.05, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.40, "tax_rate": 0.25, "capex_pct_revenue": 0.08, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.05}
    ]
  },
  {
    "id": "optimistic",
    "name": "Оптимистичный",
    "description": "Ускорение роста комиссий",
    "probability": 0.5,
    "terminal_growth_rate": 0.04,
    "assumptions": [
      {"year": 2025, "revenue_growth": 0.12, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.38, "tax_rate": 0.25, "capex_pct_revenue": 0.07, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.04},
      {"year": 2026, "revenue_growth": 0.10, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.38, "tax_rate": 0.25, "capex_pct_revenue": 0.07, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.04},
      {"year": 2027, "revenue_growth": 0.08, "cogs_pct_revenue": 0.02, "sga_pct_revenue": 0.38, "tax_rate": 0.25, "capex_pct_revenue": 0.07, "da_pct_revenue": 0.06, "nwc_pct_revenue": 0.04}
    ]
  }
]`

// TestScenarioGenerator_Execute_PriceIsNonZero проверяет, что при данных
// в формате financial-data сервиса (units=millions, акции из MOEX API)
// цена акции после DCF ненулевая.
func TestScenarioGenerator_Execute_PriceIsNonZero(t *testing.T) {
	ctx := context.Background()

	finData := mocks.NewFinancialDataGateway(t)
	aiProvider := mocks.NewAIProvider(t)
	riskRepo := mocks.NewRiskAndGrowthRepository(t)
	scenarioRepo := mocks.NewScenarioRepository(t)
	dcfRepo := mocks.NewDCFResultsRepository(t)
	transactor := mocks.NewTransactor(t)
	publisher := mocks.NewMessagePublisher(t)
	parserGateway := mocks.NewParserGateway(t)

	finData.On("GetRawDataHistory", ctx, "MOEX").Return([]entity.RawData{moexRawData()}, nil)
	finData.On("GetCBRates", ctx).Return(&entity.CBRate{Date: time.Now(), Rate: 21.0}, nil)
	finData.On("GetStockInfo", ctx, "MOEX").Return(&entity.StockInfo{
		Ticker:         "MOEX",
		NumberOfShares: 2_276_401_458,
		Name:           "Московская биржа",
	}, nil)
	// ~210 руб/акцию × 2.276 млрд акций
	finData.On("GetMarketCap", ctx, "MOEX").Return(478_044_306_180.0, nil)

	riskRepo.On("GetFreshRiskAndGrowth", ctx, "MOEX", 72*time.Hour).Return(
		&entity.RiskAndGrowthResponse{Ticker: "MOEX", Factors: []entity.RiskAndGrowthFactor{}}, nil,
	)

	aiProvider.On("GenerateText", ctx, mock.AnythingOfType("string"), entity.Pro, mock.Anything).
		Return(mockScenarioJSON, nil)

	transactor.On("RunInTx", ctx, mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		fn(ctx)
	}).Return(nil)

	scenarioRepo.On("SaveScenarios", ctx, "MOEX", mock.Anything).Return(nil)

	var capturedResult entity.DCFResult
	dcfRepo.On("SaveDCFResults", ctx, "MOEX", mock.Anything).Run(func(args mock.Arguments) {
		capturedResult = args.Get(2).(entity.DCFResult)
	}).Return(nil)

	parserGateway.On("GetLatestReport", ctx, "MOEX").Return(&entity.Report{
		Ticker: "MOEX",
		Year:   2024,
		Period: "12",
		S3Path: "s3://bucket/moex_2024_12.pdf",
	}, nil)

	publisher.On("PublishMessage", ctx, mock.Anything).Return(nil)

	sg := usecase.NewScenarioGenerator(
		aiProvider, finData, parserGateway, riskRepo,
		scenarioRepo, dcfRepo, transactor, publisher,
	)

	err := sg.Execute(ctx, entity.Task{
		Id:     "test-task-id",
		Ticker: "MOEX",
		Type:   entity.GenerateScenarios,
	})

	assert.NoError(t, err)
	assert.Len(t, capturedResult.Scenarios, 2)

	for _, s := range capturedResult.Scenarios {
		assert.Greater(t, s.PricePerShare, 0.0, "сценарий %q: цена акции ненулевая", s.ScenarioID)
		assert.Greater(t, s.EnterpriseValue, 0.0)
		assert.Greater(t, s.EquityValue, 0.0)
	}

	assert.Greater(t, capturedResult.WeightedPrice, 0.0, "взвешенная цена ненулевая")
}
