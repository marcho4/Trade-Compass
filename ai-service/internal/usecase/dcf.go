package usecase

import (
	"math"

	"ai-service/internal/domain/entity"
)

func Calculate(input entity.DCFInput, scenarios []entity.Scenario) entity.DCFResult {
	result := entity.DCFResult{}

	for _, s := range scenarios {
		sr := calculateScenario(input, s)
		result.Scenarios = append(result.Scenarios, sr)
		result.WeightedPrice += sr.PricePerShare * s.Probability
		result.WeightedEV += sr.EnterpriseValue * s.Probability
	}

	return result
}

func calculateScenario(input entity.DCFInput, scenario entity.Scenario) entity.ScenarioDCFResult {
	wacc := input.WACC
	n := len(scenario.Assumptions)

	var pvFCFs float64
	var fcfs []entity.YearlyFCF

	prevRevenue := input.BaseRevenue
	prevNWC := input.BaseNWC
	var lastFCF float64

	for i, a := range scenario.Assumptions {
		revenue := prevRevenue * (1 + a.RevenueGrowth)
		ebit := revenue * (1 - a.COGSPctRevenue - a.SGAPctRevenue)
		nopat := ebit * (1 - a.TaxRate)
		da := revenue * a.DAPctRevenue
		capex := revenue * a.CapexPctRevenue
		nwc := revenue * a.NWCPctRevenue
		deltaNWC := nwc - prevNWC

		fcf := nopat + da - capex - deltaNWC
		pvFCFs += fcf / math.Pow(1+wacc, float64(i+1))

		fcfs = append(fcfs, entity.YearlyFCF{Year: a.Year, Revenue: revenue, FCF: fcf})

		prevRevenue = revenue
		prevNWC = nwc
		lastFCF = fcf
	}

	tgr := scenario.TerminalGrowthRate
	tv := lastFCF * (1 + tgr) / (wacc - tgr)
	pvTV := tv / math.Pow(1+wacc, float64(n))

	ev := pvFCFs + pvTV
	equity := ev - input.NetDebt
	var price float64
	if input.SharesOutstanding > 0 {
		price = equity / input.SharesOutstanding
	}

	return entity.ScenarioDCFResult{
		ScenarioID:      scenario.ID,
		Probability:     scenario.Probability,
		EnterpriseValue: ev,
		EquityValue:     equity,
		PricePerShare:   price,
		TerminalValue:   tv,
		YearlyFCFs:      fcfs,
	}
}
