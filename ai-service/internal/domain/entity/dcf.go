package entity

import (
	"fmt"
	"strings"
)

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
	ID            string
	WeightedPrice float64
	WeightedEV    float64
	Scenarios     []ScenarioDCFResult
}

func (r *DCFResult) ComputeWeighted() {
	r.WeightedPrice = 0
	r.WeightedEV = 0
	for _, s := range r.Scenarios {
		r.WeightedPrice += s.PricePerShare * s.Probability
		r.WeightedEV += s.EnterpriseValue * s.Probability
	}
}

func (y *YearlyFCF) String() string {
	return fmt.Sprintf("  %d: выручка %.0f, FCF %.0f", y.Year, y.Revenue, y.FCF)
}

func (s *ScenarioDCFResult) String() string {
	var b strings.Builder

	fmt.Fprintf(&b, "--- Сценарий %s (вероятность: %.0f%%) ---\n", s.ScenarioID, s.Probability*100)
	fmt.Fprintf(&b, "Enterprise Value: %.0f\n", s.EnterpriseValue)
	fmt.Fprintf(&b, "Equity Value: %.0f\n", s.EquityValue)
	fmt.Fprintf(&b, "Цена за акцию: %.2f\n", s.PricePerShare)
	fmt.Fprintf(&b, "Терминальная стоимость: %.0f\n", s.TerminalValue)

	if len(s.YearlyFCFs) > 0 {
		b.WriteString("FCF по годам:\n")
		for _, y := range s.YearlyFCFs {
			b.WriteString(y.String())
			b.WriteByte('\n')
		}
	}

	return b.String()
}

func (r *DCFResult) String() string {
	var b strings.Builder

	b.WriteString("=== Результаты DCF-оценки ===\n")
	fmt.Fprintf(&b, "Взвешенная цена за акцию: %.2f\n", r.WeightedPrice)
	fmt.Fprintf(&b, "Взвешенный EV: %.0f\n", r.WeightedEV)

	for _, s := range r.Scenarios {
		b.WriteByte('\n')
		b.WriteString(s.String())
	}

	return b.String()
}
