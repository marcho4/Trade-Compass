package entity

import (
	"fmt"
	"strings"
)

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

func (f *RiskAndGrowthFactor) String() string {
	return fmt.Sprintf("  - [%s] %s (горизонт: %s, влияние: %s): %s | источник: %s",
		f.Type, f.Name, f.Horizon, f.Impact, f.Summary, f.Source)
}

type RiskAndGrowthResponse struct {
	Ticker  string                `json:"ticker"`
	Factors []RiskAndGrowthFactor `json:"factors"`
}

func (r *RiskAndGrowthResponse) String() string {
	var s strings.Builder

	fmt.Fprintf(&s, "=== Риски и драйверы роста: %s ===\n", r.Ticker)

	for _, f := range r.Factors {
		s.WriteString(f.String())
		s.WriteByte('\n')
	}

	return s.String()
}
