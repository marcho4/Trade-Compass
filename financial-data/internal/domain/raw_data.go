package domain

// ReportPeriod represents the reporting period for financial metrics
type ReportPeriod string

const (
	Q1   ReportPeriod = "Q1"
	Q2   ReportPeriod = "Q2"
	Q3   ReportPeriod = "Q3"
	Q4   ReportPeriod = "Q4"
	YEAR ReportPeriod = "YEAR"
)

// IsValid checks if the reporting period is valid
func (rp ReportPeriod) IsValid() bool {
	switch rp {
	case Q1, Q2, Q3, Q4, YEAR:
		return true
	default:
		return false
	}
}

// RawData contains raw financial metrics from company reports
// These metrics are used to calculate various financial ratios
type RawData struct {
	// Primary identifiers
	Ticker string       `json:"ticker"`
	Year   int          `json:"year"`
	Period ReportPeriod `json:"period"`

	// P&L (Profit & Loss Statement / Отчёт о прибылях и убытках)
	Revenue           *int64 `json:"revenue,omitempty"`            // Выручка
	CostOfRevenue     *int64 `json:"costOfRevenue,omitempty"`      // Себестоимость
	GrossProfit       *int64 `json:"grossProfit,omitempty"`        // Валовая прибыль
	OperatingExpenses *int64 `json:"operatingExpenses,omitempty"`  // Операционные расходы
	EBIT              *int64 `json:"ebit,omitempty"`               // Прибыль до вычета процентов и налогов
	EBITDA            *int64 `json:"ebitda,omitempty"`             // EBITDA
	InterestExpense   *int64 `json:"interestExpense,omitempty"`    // Проценты к уплате
	TaxExpense        *int64 `json:"taxExpense,omitempty"`         // Налоги
	NetProfit         *int64 `json:"netProfit,omitempty"`          // Чистая прибыль

	// Balance Sheet (Баланс)
	TotalAssets        *int64 `json:"totalAssets,omitempty"`        // Всего активов
	CurrentAssets      *int64 `json:"currentAssets,omitempty"`      // Оборотные активы
	CashAndEquivalents *int64 `json:"cashAndEquivalents,omitempty"` // Денежные средства и эквиваленты
	Inventories        *int64 `json:"inventories,omitempty"`        // Запасы
	Receivables        *int64 `json:"receivables,omitempty"`        // Дебиторская задолженность

	TotalLiabilities     *int64 `json:"totalLiabilities,omitempty"`     // Всего обязательств
	CurrentLiabilities   *int64 `json:"currentLiabilities,omitempty"`   // Краткосрочные обязательства
	Debt                 *int64 `json:"debt,omitempty"`                 // Долг (краткосрочный + долгосрочный)
	LongTermDebt         *int64 `json:"longTermDebt,omitempty"`         // Долгосрочный долг
	ShortTermDebt        *int64 `json:"shortTermDebt,omitempty"`        // Краткосрочный долг
	Equity               *int64 `json:"equity,omitempty"`               // Собственный капитал
	RetainedEarnings     *int64 `json:"retainedEarnings,omitempty"`     // Нераспределённая прибыль

	// Cash Flow Statement (Отчёт о движении денежных средств)
	OperatingCashFlow *int64 `json:"operatingCashFlow,omitempty"` // Операционный денежный поток
	InvestingCashFlow *int64 `json:"investingCashFlow,omitempty"` // Инвестиционный денежный поток
	FinancingCashFlow *int64 `json:"financingCashFlow,omitempty"` // Финансовый денежный поток
	CAPEX             *int64 `json:"capex,omitempty"`             // Капитальные затраты
	FreeCashFlow      *int64 `json:"freeCashFlow,omitempty"`      // Свободный денежный поток (OCF - CapEx)

	// Market Data (для мультипликаторов)
	SharesOutstanding *int64 `json:"sharesOutstanding,omitempty"` // Количество акций в обращении
	MarketCap         *int64 `json:"marketCap,omitempty"`         // Рыночная капитализация на дату отчёта

	// Calculated fields (Дополнительные расчётные поля)
	WorkingCapital  *int64 `json:"workingCapital,omitempty"`  // Оборотный капитал (current_assets - current_liabilities)
	CapitalEmployed *int64 `json:"capitalEmployed,omitempty"` // Задействованный капитал (total_assets - current_liabilities)
	EnterpriseValue *int64 `json:"enterpriseValue,omitempty"` // EV = market_cap + debt - cash
	NetDebt         *int64 `json:"netDebt,omitempty"`         // Чистый долг (debt - cash)
}
