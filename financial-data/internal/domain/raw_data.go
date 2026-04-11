package domain

type ReportPeriod string

const (
	Q1   ReportPeriod = "Q1"
	Q2   ReportPeriod = "Q2"
	Q3   ReportPeriod = "Q3"
	Q4   ReportPeriod = "Q4"
	YEAR ReportPeriod = "YEAR"
)

func (rp ReportPeriod) IsValid() bool {
	switch rp {
	case Q1, Q2, Q3, Q4, YEAR:
		return true
	default:
		return false
	}
}

type MetricsStatus string

const (
	StatusDraft     MetricsStatus = "draft"
	StatusConfirmed MetricsStatus = "confirmed"
)

func (s MetricsStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusConfirmed:
		return true
	default:
		return false
	}
}

type RawData struct {
	Ticker      string        `json:"ticker"`
	Year        int           `json:"year"`
	Period      ReportPeriod  `json:"period"`
	Status      MetricsStatus `json:"status"`
	ReportUnits *string       `json:"reportUnits,omitempty"`

	// ── Income Statement ──────────────────────────────────────────
	Revenue           *int64   `json:"revenue,omitempty"`
	CostOfRevenue     *int64   `json:"costOfRevenue,omitempty"`
	GrossProfit       *int64   `json:"grossProfit,omitempty"`
	OperatingExpenses *int64   `json:"operatingExpenses,omitempty"`
	OtherIncome       *int64   `json:"otherIncome,omitempty"`
	OtherExpenses     *int64   `json:"otherExpenses,omitempty"`
	EBIT              *int64   `json:"ebit,omitempty"`
	EBITDA            *int64   `json:"ebitda,omitempty"`
	Depreciation      *int64   `json:"depreciation,omitempty"`
	InterestIncome    *int64   `json:"interestIncome,omitempty"`
	InterestExpense   *int64   `json:"interestExpense,omitempty"`
	ProfitBeforeTax   *int64   `json:"profitBeforeTax,omitempty"`
	TaxExpense        *int64   `json:"taxExpense,omitempty"`
	NetProfit         *int64   `json:"netProfit,omitempty"`
	NetProfitParent   *int64   `json:"netProfitParent,omitempty"`
	BasicEPS          *float64 `json:"basicEps,omitempty"`

	// ── Balance Sheet ─────────────────────────────────────────────
	TotalAssets           *int64 `json:"totalAssets,omitempty"`
	CurrentAssets         *int64 `json:"currentAssets,omitempty"`
	CashAndEquivalents    *int64 `json:"cashAndEquivalents,omitempty"`
	Inventories           *int64 `json:"inventories,omitempty"`
	Receivables           *int64 `json:"receivables,omitempty"`
	FixedAssets           *int64 `json:"fixedAssets,omitempty"`
	RightOfUseAssets      *int64 `json:"rightOfUseAssets,omitempty"`
	IntangibleAssets      *int64 `json:"intangibleAssets,omitempty"`
	Goodwill              *int64 `json:"goodwill,omitempty"`
	TotalNonCurrentAssets *int64 `json:"totalNonCurrentAssets,omitempty"`

	TotalLiabilities   *int64 `json:"totalLiabilities,omitempty"`
	CurrentLiabilities *int64 `json:"currentLiabilities,omitempty"`
	Debt               *int64 `json:"debt,omitempty"`
	LongTermDebt       *int64 `json:"longTermDebt,omitempty"`
	ShortTermDebt      *int64 `json:"shortTermDebt,omitempty"`
	LtLeaseLiabilities *int64 `json:"ltLeaseLiabilities,omitempty"`
	StLeaseLiabilities *int64 `json:"stLeaseLiabilities,omitempty"`
	TradePayables      *int64 `json:"tradePayables,omitempty"`
	Equity             *int64 `json:"equity,omitempty"`
	EquityParent       *int64 `json:"equityParent,omitempty"`
	TreasuryShares     *int64 `json:"treasuryShares,omitempty"`
	RetainedEarnings   *int64 `json:"retainedEarnings,omitempty"`

	// ── Cash Flow ─────────────────────────────────────────────────
	OperatingCashFlow *int64 `json:"operatingCashFlow,omitempty"`
	InvestingCashFlow *int64 `json:"investingCashFlow,omitempty"`
	FinancingCashFlow *int64 `json:"financingCashFlow,omitempty"`
	CAPEX             *int64 `json:"capex,omitempty"`
	FreeCashFlow      *int64 `json:"freeCashFlow,omitempty"`
	DividendsPaid     *int64 `json:"dividendsPaid,omitempty"`
	LeasePayments     *int64 `json:"leasePayments,omitempty"`
	AcquisitionsNet   *int64 `json:"acquisitionsNet,omitempty"`
	InterestPaid      *int64 `json:"interestPaid,omitempty"`
	DebtProceeds      *int64 `json:"debtProceeds,omitempty"`
	DebtRepayments    *int64 `json:"debtRepayments,omitempty"`

	// ── Per Share & Market ────────────────────────────────────────
	SharesOutstanding *int64 `json:"sharesOutstanding,omitempty"`
	MarketCap         *int64 `json:"marketCap,omitempty"`
	EnterpriseValue   *int64 `json:"enterpriseValue,omitempty"`

	// ── Derived Metrics ───────────────────────────────────────────
	WorkingCapital  *int64 `json:"workingCapital,omitempty"`
	CapitalEmployed *int64 `json:"capitalEmployed,omitempty"`
	NetDebt         *int64 `json:"netDebt,omitempty"`

	// ── Notes Breakdown ───────────────────────────────────────────
	InterestOnLeases *int64 `json:"interestOnLeases,omitempty"`
	InterestOnLoans  *int64 `json:"interestOnLoans,omitempty"`

	// ── Bank-specific ─────────────────────────────────────────────
	CompanyType          *string `json:"companyType,omitempty"`
	NetInterestIncome    *int64  `json:"netInterestIncome,omitempty"`
	CommissionIncome     *int64  `json:"commissionIncome,omitempty"`
	CommissionExpense    *int64  `json:"commissionExpense,omitempty"`
	NetCommissionIncome  *int64  `json:"netCommissionIncome,omitempty"`
	CreditLossProvision  *int64  `json:"creditLossProvision,omitempty"`
}
