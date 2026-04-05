package gemini

import "ai-service/internal/domain/entity"

type scenarioDTO struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	Description          string               `json:"description"`
	Probability          float64              `json:"probability"`
	TerminalGrowthRate   float64              `json:"terminal_growth_rate"`
	GrowthFactorsApplied []factorDTO          `json:"growth_factors_applied,omitempty"`
	RisksApplied         []factorDTO          `json:"risks_applied,omitempty"`
	Assumptions          []yearlyAssumptionDTO `json:"assumptions"`
}

type factorDTO struct {
	Factor string `json:"factor"`
	Impact string `json:"impact"`
}

type yearlyAssumptionDTO struct {
	Year            int     `json:"year"`
	RevenueGrowth   float64 `json:"revenue_growth"`
	COGSPctRevenue  float64 `json:"cogs_pct_revenue"`
	SGAPctRevenue   float64 `json:"sga_pct_revenue"`
	TaxRate         float64 `json:"tax_rate"`
	CapexPctRevenue float64 `json:"capex_pct_revenue"`
	DAPctRevenue    float64 `json:"da_pct_revenue"`
	NWCPctRevenue   float64 `json:"nwc_pct_revenue"`
}

func mapScenariosToDomain(dtos []scenarioDTO) []entity.Scenario {
	scenarios := make([]entity.Scenario, len(dtos))
	for i, d := range dtos {
		scenarios[i] = entity.Scenario{
			ID:                 d.ID,
			Name:               d.Name,
			Description:        d.Description,
			Probability:        d.Probability,
			TerminalGrowthRate: d.TerminalGrowthRate,
			GrowthFactorsApplied: mapFactorsToDomain(d.GrowthFactorsApplied),
			RisksApplied:         mapFactorsToDomain(d.RisksApplied),
			Assumptions:          mapAssumptionsToDomain(d.Assumptions),
		}
	}
	return scenarios
}

func mapFactorsToDomain(dtos []factorDTO) []entity.Factor {
	if dtos == nil {
		return nil
	}
	factors := make([]entity.Factor, len(dtos))
	for i, d := range dtos {
		factors[i] = entity.Factor{
			Factor: d.Factor,
			Impact: d.Impact,
		}
	}
	return factors
}

func mapAssumptionsToDomain(dtos []yearlyAssumptionDTO) []entity.YearlyAssumption {
	assumptions := make([]entity.YearlyAssumption, len(dtos))
	for i, d := range dtos {
		assumptions[i] = entity.YearlyAssumption{
			Year:            d.Year,
			RevenueGrowth:   d.RevenueGrowth,
			COGSPctRevenue:  d.COGSPctRevenue,
			SGAPctRevenue:   d.SGAPctRevenue,
			TaxRate:         d.TaxRate,
			CapexPctRevenue: d.CapexPctRevenue,
			DAPctRevenue:    d.DAPctRevenue,
			NWCPctRevenue:   d.NWCPctRevenue,
		}
	}
	return assumptions
}
