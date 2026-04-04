package entity

type DCFInput struct {
	BaseRevenue       float64
	BaseNWC           float64
	WACC              float64
	NetDebt           float64
	SharesOutstanding int64
}

type YearlyFCF struct {
	Year    int
	Revenue float64
	FCF     float64
}

type ScenarioDCFResult struct {
	ScenarioID      string
	Probability     float64
	EnterpriseValue float64
	EquityValue     float64
	PricePerShare   float64
	TerminalValue   float64
	YearlyFCFs      []YearlyFCF
}

type DCFResult struct {
	WeightedPrice float64
	WeightedEV    float64
	Scenarios     []ScenarioDCFResult
}
