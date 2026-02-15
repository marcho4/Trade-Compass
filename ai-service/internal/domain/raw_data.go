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
	Ticker string       `json:"ticker"`
	Year   int          `json:"year"`
	Period ReportPeriod `json:"period"`
	Status string       `json:"status"`

	Revenue           *int64 `json:"revenue,omitempty"`
	CostOfRevenue     *int64 `json:"costOfRevenue,omitempty"`
	GrossProfit       *int64 `json:"grossProfit,omitempty"`
	OperatingExpenses *int64 `json:"operatingExpenses,omitempty"`
	EBIT              *int64 `json:"ebit,omitempty"`
	EBITDA            *int64 `json:"ebitda,omitempty"`
	InterestExpense   *int64 `json:"interestExpense,omitempty"`
	TaxExpense        *int64 `json:"taxExpense,omitempty"`
	NetProfit         *int64 `json:"netProfit,omitempty"`

	TotalAssets        *int64 `json:"totalAssets,omitempty"`
	CurrentAssets      *int64 `json:"currentAssets,omitempty"`
	CashAndEquivalents *int64 `json:"cashAndEquivalents,omitempty"`
	Inventories        *int64 `json:"inventories,omitempty"`
	Receivables        *int64 `json:"receivables,omitempty"`

	TotalLiabilities   *int64 `json:"totalLiabilities,omitempty"`
	CurrentLiabilities *int64 `json:"currentLiabilities,omitempty"`
	Debt               *int64 `json:"debt,omitempty"`
	LongTermDebt       *int64 `json:"longTermDebt,omitempty"`
	ShortTermDebt      *int64 `json:"shortTermDebt,omitempty"`
	Equity             *int64 `json:"equity,omitempty"`
	RetainedEarnings   *int64 `json:"retainedEarnings,omitempty"`

	OperatingCashFlow *int64 `json:"operatingCashFlow,omitempty"`
	InvestingCashFlow *int64 `json:"investingCashFlow,omitempty"`
	FinancingCashFlow *int64 `json:"financingCashFlow,omitempty"`
	CAPEX             *int64 `json:"capex,omitempty"`
	FreeCashFlow      *int64 `json:"freeCashFlow,omitempty"`

	SharesOutstanding *int64 `json:"sharesOutstanding,omitempty"`
	MarketCap         *int64 `json:"marketCap,omitempty"`

	WorkingCapital  *int64 `json:"workingCapital,omitempty"`
	CapitalEmployed *int64 `json:"capitalEmployed,omitempty"`
	EnterpriseValue *int64 `json:"enterpriseValue,omitempty"`
	NetDebt         *int64 `json:"netDebt,omitempty"`
}

type Report struct {
	ID     int    `json:"id"`
	Ticker string `json:"ticker"`
	Year   int    `json:"year"`
	Period string `json:"period"`
	S3Path string `json:"s3_path"`
}
