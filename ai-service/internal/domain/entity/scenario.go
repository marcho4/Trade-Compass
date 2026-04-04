package entity

type ScenariosResponse struct {
	Scenarios []Scenario `json:"scenarios"`
}

type Scenario struct {
	ID                   string             `json:"id"`
	Name                 string             `json:"name"`
	Description          string             `json:"description"`
	Probability          float64            `json:"probability"`
	TerminalGrowthRate   float64            `json:"terminal_growth_rate"`
	GrowthFactorsApplied []Factor           `json:"growth_factors_applied,omitempty"`
	RisksApplied         []Factor           `json:"risks_applied,omitempty"`
	Assumptions          []YearlyAssumption `json:"assumptions"`
}

type Factor struct {
	Factor string `json:"factor"`
	Impact string `json:"impact"`
}

type YearlyAssumption struct {
	Year            int     `json:"year"`
	RevenueGrowth   float64 `json:"revenue_growth"`
	COGSPctRevenue  float64 `json:"cogs_pct_revenue"`
	SGAPctRevenue   float64 `json:"sga_pct_revenue"`
	TaxRate         float64 `json:"tax_rate"`
	CapexPctRevenue float64 `json:"capex_pct_revenue"`
	DAPctRevenue    float64 `json:"da_pct_revenue"`
	NWCPctRevenue   float64 `json:"nwc_pct_revenue"`
}
