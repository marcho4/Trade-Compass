package domain

import "sort"

func reportUnitsMultiplier(units *string) int64 {
	if units == nil {
		return 1
	}
	switch *units {
	case "thousands":
		return 1_000
	case "millions":
		return 1_000_000
	case "billions":
		return 1_000_000_000
	case "trillions":
		return 1_000_000_000_000
	default:
		return 1
	}
}

func addScaledPtr(dst **int64, src *int64, mul int64) {
	if src == nil {
		return
	}
	if *dst == nil {
		v := int64(0)
		*dst = &v
	}
	**dst += *src * mul
}

func accumulateRawData(dst, src *RawData, sign int64) {
	mul := sign * reportUnitsMultiplier(src.ReportUnits)

	addScaledPtr(&dst.Revenue, src.Revenue, mul)
	addScaledPtr(&dst.CostOfRevenue, src.CostOfRevenue, mul)
	addScaledPtr(&dst.GrossProfit, src.GrossProfit, mul)
	addScaledPtr(&dst.OperatingExpenses, src.OperatingExpenses, mul)
	addScaledPtr(&dst.OtherIncome, src.OtherIncome, mul)
	addScaledPtr(&dst.OtherExpenses, src.OtherExpenses, mul)
	addScaledPtr(&dst.EBIT, src.EBIT, mul)
	addScaledPtr(&dst.EBITDA, src.EBITDA, mul)
	addScaledPtr(&dst.Depreciation, src.Depreciation, mul)
	addScaledPtr(&dst.InterestIncome, src.InterestIncome, mul)
	addScaledPtr(&dst.InterestExpense, src.InterestExpense, mul)
	addScaledPtr(&dst.ProfitBeforeTax, src.ProfitBeforeTax, mul)
	addScaledPtr(&dst.TaxExpense, src.TaxExpense, mul)
	addScaledPtr(&dst.NetProfit, src.NetProfit, mul)
	addScaledPtr(&dst.NetProfitParent, src.NetProfitParent, mul)

	addScaledPtr(&dst.TotalAssets, src.TotalAssets, mul)
	addScaledPtr(&dst.CurrentAssets, src.CurrentAssets, mul)
	addScaledPtr(&dst.CashAndEquivalents, src.CashAndEquivalents, mul)
	addScaledPtr(&dst.Inventories, src.Inventories, mul)
	addScaledPtr(&dst.Receivables, src.Receivables, mul)
	addScaledPtr(&dst.FixedAssets, src.FixedAssets, mul)
	addScaledPtr(&dst.RightOfUseAssets, src.RightOfUseAssets, mul)
	addScaledPtr(&dst.IntangibleAssets, src.IntangibleAssets, mul)
	addScaledPtr(&dst.Goodwill, src.Goodwill, mul)
	addScaledPtr(&dst.TotalNonCurrentAssets, src.TotalNonCurrentAssets, mul)

	addScaledPtr(&dst.TotalLiabilities, src.TotalLiabilities, mul)
	addScaledPtr(&dst.CurrentLiabilities, src.CurrentLiabilities, mul)
	addScaledPtr(&dst.Debt, src.Debt, mul)
	addScaledPtr(&dst.LongTermDebt, src.LongTermDebt, mul)
	addScaledPtr(&dst.ShortTermDebt, src.ShortTermDebt, mul)
	addScaledPtr(&dst.LtLeaseLiabilities, src.LtLeaseLiabilities, mul)
	addScaledPtr(&dst.StLeaseLiabilities, src.StLeaseLiabilities, mul)
	addScaledPtr(&dst.TradePayables, src.TradePayables, mul)
	addScaledPtr(&dst.Equity, src.Equity, mul)
	addScaledPtr(&dst.EquityParent, src.EquityParent, mul)
	addScaledPtr(&dst.TreasuryShares, src.TreasuryShares, mul)
	addScaledPtr(&dst.RetainedEarnings, src.RetainedEarnings, mul)

	addScaledPtr(&dst.OperatingCashFlow, src.OperatingCashFlow, mul)
	addScaledPtr(&dst.InvestingCashFlow, src.InvestingCashFlow, mul)
	addScaledPtr(&dst.FinancingCashFlow, src.FinancingCashFlow, mul)
	addScaledPtr(&dst.CAPEX, src.CAPEX, mul)
	addScaledPtr(&dst.FreeCashFlow, src.FreeCashFlow, mul)
	addScaledPtr(&dst.DividendsPaid, src.DividendsPaid, mul)
	addScaledPtr(&dst.LeasePayments, src.LeasePayments, mul)
	addScaledPtr(&dst.AcquisitionsNet, src.AcquisitionsNet, mul)
	addScaledPtr(&dst.InterestPaid, src.InterestPaid, mul)
	addScaledPtr(&dst.DebtProceeds, src.DebtProceeds, mul)
	addScaledPtr(&dst.DebtRepayments, src.DebtRepayments, mul)

	addScaledPtr(&dst.SharesOutstanding, src.SharesOutstanding, mul)
	addScaledPtr(&dst.MarketCap, src.MarketCap, mul)
	addScaledPtr(&dst.EnterpriseValue, src.EnterpriseValue, mul)

	addScaledPtr(&dst.WorkingCapital, src.WorkingCapital, mul)
	addScaledPtr(&dst.CapitalEmployed, src.CapitalEmployed, mul)
	addScaledPtr(&dst.NetDebt, src.NetDebt, mul)

	addScaledPtr(&dst.InterestOnLeases, src.InterestOnLeases, mul)
	addScaledPtr(&dst.InterestOnLoans, src.InterestOnLoans, mul)

	addScaledPtr(&dst.NetInterestIncome, src.NetInterestIncome, mul)
	addScaledPtr(&dst.CommissionIncome, src.CommissionIncome, mul)
	addScaledPtr(&dst.CommissionExpense, src.CommissionExpense, mul)
	addScaledPtr(&dst.NetCommissionIncome, src.NetCommissionIncome, mul)
	addScaledPtr(&dst.CreditLossProvision, src.CreditLossProvision, mul)
}

func BuildAnnualSnapshots(history []RawData) []RawData {
	byYear := make(map[int]map[ReportPeriod]*RawData)
	for i := range history {
		row := &history[i]
		if _, ok := byYear[row.Year]; !ok {
			byYear[row.Year] = make(map[ReportPeriod]*RawData)
		}
		byYear[row.Year][row.Period] = row
	}

	years := make([]int, 0, len(byYear))
	for y := range byYear {
		years = append(years, y)
	}
	sort.Ints(years)

	type keyed struct {
		sortKey int
		snap    RawData
	}
	var result []keyed

	for _, y := range years {
		periods := byYear[y]
		prev := byYear[y-1]

		if yr, ok := periods[YEAR]; ok {
			snap := RawData{Ticker: yr.Ticker, Year: y, Period: YEAR, Status: yr.Status}
			accumulateRawData(&snap, yr, 1)
			result = append(result, keyed{sortKey: y*100 + 12, snap: snap})
		}

		if q2, ok := periods[Q2]; ok && prev != nil {
			if prevYear, okY := prev[YEAR]; okY {
				if prevQ2, okQ := prev[Q2]; okQ {
					snap := RawData{Ticker: q2.Ticker, Year: y, Period: Q2, Status: q2.Status}
					accumulateRawData(&snap, prevYear, 1)
					accumulateRawData(&snap, prevQ2, -1)
					accumulateRawData(&snap, q2, 1)
					result = append(result, keyed{sortKey: y*100 + 6, snap: snap})
				}
			}
		}

		if q3, ok := periods[Q3]; ok && prev != nil {
			if prevYear, okY := prev[YEAR]; okY {
				if prevQ3, okQ := prev[Q3]; okQ {
					snap := RawData{Ticker: q3.Ticker, Year: y, Period: Q3, Status: q3.Status}
					accumulateRawData(&snap, prevYear, 1)
					accumulateRawData(&snap, prevQ3, -1)
					accumulateRawData(&snap, q3, 1)
					result = append(result, keyed{sortKey: y*100 + 9, snap: snap})
				}
			}
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].sortKey < result[j].sortKey
	})

	out := make([]RawData, len(result))
	for i, k := range result {
		out[i] = k.snap
	}
	return out
}
