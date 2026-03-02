package domain

type ReportPeriod string

const (
	Q1   ReportPeriod = "Q1"
	Q2   ReportPeriod = "Q2"
	Q3   ReportPeriod = "Q3"
	Q4   ReportPeriod = "Q4"
	YEAR ReportPeriod = "YEAR"
)

var MonthsToPeriod = map[string]ReportPeriod{
	"3":  Q1,
	"6":  Q2,
	"9":  Q3,
	"12": YEAR,
}

var PeriodToMonths = map[string]int{
	string(Q1):   3,
	string(Q2):   6,
	string(Q3):   9,
	string(YEAR): 12,
}

type RawData struct {
	Ticker      string       `json:"ticker"`
	Year        int          `json:"year"`
	Period      ReportPeriod `json:"period"`
	Status      string       `json:"status"`
	ReportUnits string       `json:"reportUnits"` // "thousands" | "millions" | "units"

	// ── Income Statement ──────────────────────────────────────────
	Revenue           *int64   `json:"revenue,omitempty"`
	CostOfRevenue     *int64   `json:"costOfRevenue,omitempty"`
	GrossProfit       *int64   `json:"grossProfit,omitempty"`
	OperatingExpenses *int64   `json:"operatingExpenses,omitempty"` // SG&A
	OtherIncome       *int64   `json:"otherIncome,omitempty"`       // Прочие доходы + доход от аренды
	OtherExpenses     *int64   `json:"otherExpenses,omitempty"`     // Прочие расходы (отрицательное)
	EBIT              *int64   `json:"ebit,omitempty"`
	EBITDA            *int64   `json:"ebitda,omitempty"`          // расчётное: EBIT + Depreciation
	Depreciation      *int64   `json:"depreciation,omitempty"`    // Общая D&A из CF
	InterestIncome    *int64   `json:"interestIncome,omitempty"`  // Процентные доходы
	InterestExpense   *int64   `json:"interestExpense,omitempty"` // Финансовые расходы
	ProfitBeforeTax   *int64   `json:"profitBeforeTax,omitempty"`
	TaxExpense        *int64   `json:"taxExpense,omitempty"`
	NetProfit         *int64   `json:"netProfit,omitempty"`
	NetProfitParent   *int64   `json:"netProfitParent,omitempty"` // ЧП на акционеров материнской
	BasicEPS          *float64 `json:"basicEps,omitempty"`        // Базовая прибыль на акцию (руб.)

	// ── Balance Sheet ─────────────────────────────────────────────
	TotalAssets           *int64 `json:"totalAssets,omitempty"`
	CurrentAssets         *int64 `json:"currentAssets,omitempty"`
	CashAndEquivalents    *int64 `json:"cashAndEquivalents,omitempty"`
	Inventories           *int64 `json:"inventories,omitempty"`
	Receivables           *int64 `json:"receivables,omitempty"`
	FixedAssets           *int64 `json:"fixedAssets,omitempty"`      // Основные средства
	RightOfUseAssets      *int64 `json:"rightOfUseAssets,omitempty"` // Активы ППА (IFRS 16)
	IntangibleAssets      *int64 `json:"intangibleAssets,omitempty"` // НМА
	Goodwill              *int64 `json:"goodwill,omitempty"`
	TotalNonCurrentAssets *int64 `json:"totalNonCurrentAssets,omitempty"`

	TotalLiabilities   *int64 `json:"totalLiabilities,omitempty"`
	CurrentLiabilities *int64 `json:"currentLiabilities,omitempty"`
	Debt               *int64 `json:"debt,omitempty"` // расчётное: LT + ST
	LongTermDebt       *int64 `json:"longTermDebt,omitempty"`
	ShortTermDebt      *int64 `json:"shortTermDebt,omitempty"`
	LtLeaseLiabilities *int64 `json:"ltLeaseLiabilities,omitempty"` // Долгосрочные обязательства по аренде
	StLeaseLiabilities *int64 `json:"stLeaseLiabilities,omitempty"` // Краткосрочные обязательства по аренде
	TradePayables      *int64 `json:"tradePayables,omitempty"`      // Торговая кредиторка
	Equity             *int64 `json:"equity,omitempty"`             // Итого капитал (с НКД)
	EquityParent       *int64 `json:"equityParent,omitempty"`       // Капитал акционеров материнской
	TreasuryShares     *int64 `json:"treasuryShares,omitempty"`     // Казначейские акции (отрицательное)
	RetainedEarnings   *int64 `json:"retainedEarnings,omitempty"`

	// ── Cash Flow ─────────────────────────────────────────────────
	OperatingCashFlow *int64 `json:"operatingCashFlow,omitempty"`
	InvestingCashFlow *int64 `json:"investingCashFlow,omitempty"`
	FinancingCashFlow *int64 `json:"financingCashFlow,omitempty"`
	CAPEX             *int64 `json:"capex,omitempty"`           // расчётное: capexFA + capexIA (отрицательное)
	FreeCashFlow      *int64 `json:"freeCashFlow,omitempty"`    // расчётное: OCF + CAPEX
	DividendsPaid     *int64 `json:"dividendsPaid,omitempty"`   // Дивиденды выплаченные (отрицательное)
	LeasePayments     *int64 `json:"leasePayments,omitempty"`   // Погашение обязательств по аренде (отрицательное)
	AcquisitionsNet   *int64 `json:"acquisitionsNet,omitempty"` // Покупки бизнесов нетто (отрицательное)
	InterestPaid      *int64 `json:"interestPaid,omitempty"`    // Проценты уплаченные (отрицательное)
	DebtProceeds      *int64 `json:"debtProceeds,omitempty"`    // Привлечение кредитов
	DebtRepayments    *int64 `json:"debtRepayments,omitempty"`  // Погашение кредитов (отрицательное)

	// ── Per Share & Market ────────────────────────────────────────
	SharesOutstanding *int64 `json:"sharesOutstanding,omitempty"` // Акции в обращении (штуки)
	MarketCap         *int64 `json:"marketCap,omitempty"`         // Заполняется из MOEX API
	EnterpriseValue   *int64 `json:"enterpriseValue,omitempty"`   // Заполняется из MOEX API

	// ── Derived Metrics ───────────────────────────────────────────
	WorkingCapital  *int64 `json:"workingCapital,omitempty"`  // currentAssets - currentLiabilities
	CapitalEmployed *int64 `json:"capitalEmployed,omitempty"` // totalAssets - currentLiabilities
	NetDebt         *int64 `json:"netDebt,omitempty"`         // debt - cash (без аренды)

	// ── Notes Breakdown ───────────────────────────────────────────
	InterestOnLeases *int64 `json:"interestOnLeases,omitempty"` // Проценты по аренде (IFRS 16)
	InterestOnLoans  *int64 `json:"interestOnLoans,omitempty"`  // Проценты по кредитам + облигациям

	// ── Validation ────────────────────────────────────────────────
	Warnings []string `json:"warnings,omitempty"`
}

type Report struct {
	ID     int    `json:"id"`
	Ticker string `json:"ticker"`
	Year   int    `json:"year"`
	Period string `json:"period"`
	S3Path string `json:"s3_path"`
}
