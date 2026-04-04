package entity

type FactorType string

const (
	FactorGrowth FactorType = "growth"
	FactorRisk   FactorType = "risk"
)

type FactorHorizon string

const (
	HorizonShortTerm  FactorHorizon = "short_term"
	HorizonMediumTerm FactorHorizon = "medium_term"
)

type FactorImpact string

const (
	ImpactHigh   FactorImpact = "high"
	ImpactMedium FactorImpact = "medium"
	ImpactLow    FactorImpact = "low"
)

type RiskAndGrowthFactor struct {
	Name    string        `json:"name"`
	Type    FactorType    `json:"type"`
	Horizon FactorHorizon `json:"horizon"`
	Impact  FactorImpact  `json:"impact"`
	Summary string        `json:"summary"`
	Source  string        `json:"source"`
}

type RiskAndGrowthResponse struct {
	Ticker  string                `json:"ticker"`
	Factors []RiskAndGrowthFactor `json:"factors"`
}
