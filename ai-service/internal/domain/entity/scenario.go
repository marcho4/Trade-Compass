package entity

type Scenario struct {
	ID                   string
	Name                 string
	Description          string
	Probability          float64
	TerminalGrowthRate   float64
	GrowthFactorsApplied []Factor
	RisksApplied         []Factor
	Assumptions          []YearlyAssumption
}

type Factor struct {
	Factor string
	Impact string
}

type YearlyAssumption struct {
	Year            int
	RevenueGrowth   float64
	COGSPctRevenue  float64
	SGAPctRevenue   float64
	TaxRate         float64
	CapexPctRevenue float64
	DAPctRevenue    float64
	NWCPctRevenue   float64
}
