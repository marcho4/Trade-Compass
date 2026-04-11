package entity

import (
	"fmt"
	"strings"
)

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

func (f *Factor) String() string {
	return fmt.Sprintf("  - %s (влияние: %s)", f.Factor, f.Impact)
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

func (a *YearlyAssumption) String() string {
	return fmt.Sprintf("  %d: рост выручки %.1f%%, COGS %.1f%%, SGA %.1f%%, налог %.1f%%, CAPEX %.1f%%, DA %.1f%%, NWC %.1f%%",
		a.Year, a.RevenueGrowth*100, a.COGSPctRevenue*100, a.SGAPctRevenue*100, a.TaxRate*100, a.CapexPctRevenue*100, a.DAPctRevenue*100, a.NWCPctRevenue*100)
}

func (s *Scenario) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "--- Сценарий: %s (вероятность: %.0f%%) ---\n", s.Name, s.Probability*100)
	fmt.Fprintf(&b, "Описание: %s\n", s.Description)
	fmt.Fprintf(&b, "Терминальный темп роста: %.2f%%\n", s.TerminalGrowthRate*100)

	if len(s.GrowthFactorsApplied) > 0 {
		b.WriteString("Драйверы роста:\n")
		for _, f := range s.GrowthFactorsApplied {
			b.WriteString(f.String())
			b.WriteByte('\n')
		}
	}

	if len(s.RisksApplied) > 0 {
		b.WriteString("Риски:\n")
		for _, f := range s.RisksApplied {
			b.WriteString(f.String())
			b.WriteByte('\n')
		}
	}

	if len(s.Assumptions) > 0 {
		b.WriteString("Допущения по годам:\n")
		for _, a := range s.Assumptions {
			b.WriteString(a.String())
			b.WriteByte('\n')
		}
	}

	return b.String()
}
