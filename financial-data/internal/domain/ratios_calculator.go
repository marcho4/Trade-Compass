package domain

import "math"

func ptr(v float64) *float64 {
	return &v
}

func safeDiv(a, b float64) *float64 {
	if b == 0 {
		return nil
	}
	return ptr(a / b)
}

func pctGrowth(current, previous float64) *float64 {
	if previous == 0 {
		return nil
	}
	return ptr((current - previous) / math.Abs(previous) * 100)
}

func CalculateRatios(current *RawData, previous *RawData) *Ratios {
	r := &Ratios{}

	var marketCap float64
	hasMarketCap := current.MarketCap != nil
	if hasMarketCap {
		marketCap = float64(*current.MarketCap)
		r.MarketCap = ptr(marketCap)
	}

	var ev float64
	if hasMarketCap && current.NetDebt != nil {
		ev = marketCap + float64(*current.NetDebt)
		r.EnterpriseValue = ptr(ev)
	}

	if hasMarketCap && current.NetProfitParent != nil {
		r.PriceToEarnings = safeDiv(marketCap, float64(*current.NetProfitParent))
	}

	if hasMarketCap && current.EquityParent != nil {
		r.PriceToBook = safeDiv(marketCap, float64(*current.EquityParent))
	}

	if hasMarketCap && current.OperatingCashFlow != nil {
		r.PriceToCashFlow = safeDiv(marketCap, float64(*current.OperatingCashFlow))
	}

	if r.EnterpriseValue != nil && current.EBITDA != nil {
		r.EVToEBITDA = safeDiv(ev, float64(*current.EBITDA))
	}

	if r.EnterpriseValue != nil && current.Revenue != nil {
		r.EVToSales = safeDiv(ev, float64(*current.Revenue))
	}

	if r.EnterpriseValue != nil && current.FreeCashFlow != nil {
		r.EVToFCF = safeDiv(ev, float64(*current.FreeCashFlow))
	}

	if current.NetProfit != nil && current.Equity != nil {
		r.ROE = safeDiv(float64(*current.NetProfit), float64(*current.Equity))
		if r.ROE != nil {
			r.ROE = ptr(*r.ROE * 100)
		}
	}

	if current.NetProfit != nil && current.TotalAssets != nil {
		r.ROA = safeDiv(float64(*current.NetProfit), float64(*current.TotalAssets))
		if r.ROA != nil {
			r.ROA = ptr(*r.ROA * 100)
		}
	}

	if current.EBIT != nil && current.TaxExpense != nil && current.ProfitBeforeTax != nil && *current.ProfitBeforeTax != 0 && current.CapitalEmployed != nil {
		taxRate := float64(*current.TaxExpense) / float64(*current.ProfitBeforeTax)
		nopat := float64(*current.EBIT) * (1 - taxRate)
		r.ROIC = safeDiv(nopat, float64(*current.CapitalEmployed))
		if r.ROIC != nil {
			r.ROIC = ptr(*r.ROIC * 100)
		}
	}

	if current.GrossProfit != nil && current.Revenue != nil {
		r.GrossProfitMargin = safeDiv(float64(*current.GrossProfit), float64(*current.Revenue))
		if r.GrossProfitMargin != nil {
			r.GrossProfitMargin = ptr(*r.GrossProfitMargin * 100)
		}
	}

	if current.EBIT != nil && current.Revenue != nil {
		r.OperatingProfitMargin = safeDiv(float64(*current.EBIT), float64(*current.Revenue))
		if r.OperatingProfitMargin != nil {
			r.OperatingProfitMargin = ptr(*r.OperatingProfitMargin * 100)
		}
	}

	if current.NetProfit != nil && current.Revenue != nil {
		r.NetProfitMargin = safeDiv(float64(*current.NetProfit), float64(*current.Revenue))
		if r.NetProfitMargin != nil {
			r.NetProfitMargin = ptr(*r.NetProfitMargin * 100)
		}
	}

	if current.CurrentAssets != nil && current.CurrentLiabilities != nil {
		r.CurrentRatio = safeDiv(float64(*current.CurrentAssets), float64(*current.CurrentLiabilities))
	}

	if current.CashAndEquivalents != nil && current.Receivables != nil && current.CurrentLiabilities != nil {
		quick := float64(*current.CashAndEquivalents) + float64(*current.Receivables)
		r.QuickRatio = safeDiv(quick, float64(*current.CurrentLiabilities))
	}

	if current.NetDebt != nil && current.EBITDA != nil {
		r.NetDebtToEBITDA = safeDiv(float64(*current.NetDebt), float64(*current.EBITDA))
	}

	if current.Debt != nil && current.Equity != nil {
		r.DebtToEquity = safeDiv(float64(*current.Debt), float64(*current.Equity))
	}

	if current.EBIT != nil && current.InterestExpense != nil {
		r.InterestCoverageRatio = safeDiv(float64(*current.EBIT), math.Abs(float64(*current.InterestExpense)))
	}

	if current.OperatingCashFlow != nil && current.NetProfit != nil {
		r.IncomeQuality = safeDiv(float64(*current.OperatingCashFlow), float64(*current.NetProfit))
	}

	if current.Revenue != nil && current.TotalAssets != nil {
		r.AssetTurnover = safeDiv(float64(*current.Revenue), float64(*current.TotalAssets))
	}

	if current.CostOfRevenue != nil && current.Inventories != nil {
		r.InventoryTurnover = safeDiv(float64(*current.CostOfRevenue), float64(*current.Inventories))
	}

	if current.Revenue != nil && current.Receivables != nil {
		r.ReceivablesTurnover = safeDiv(float64(*current.Revenue), float64(*current.Receivables))
	}

	if current.BasicEPS != nil {
		r.EPS = ptr(*current.BasicEPS)
	}

	hasShares := current.SharesOutstanding != nil
	var shares float64
	if hasShares {
		shares = float64(*current.SharesOutstanding)
	}

	if hasShares && current.EquityParent != nil {
		r.BookValuePerShare = safeDiv(float64(*current.EquityParent), shares)
	}

	if hasShares && current.OperatingCashFlow != nil {
		r.CashFlowPerShare = safeDiv(float64(*current.OperatingCashFlow), shares)
	}

	var dividendPerShare *float64
	if hasShares && current.DividendsPaid != nil {
		dividendPerShare = safeDiv(math.Abs(float64(*current.DividendsPaid)), shares)
		r.DividendPerShare = dividendPerShare
	}

	if dividendPerShare != nil && hasMarketCap && hasShares {
		pricePerShare := marketCap / shares
		r.DividendYield = safeDiv(*dividendPerShare, pricePerShare)
		if r.DividendYield != nil {
			r.DividendYield = ptr(*r.DividendYield * 100)
		}
	}

	if current.DividendsPaid != nil && current.NetProfit != nil {
		r.PayoutRatio = safeDiv(math.Abs(float64(*current.DividendsPaid)), float64(*current.NetProfit))
		if r.PayoutRatio != nil {
			r.PayoutRatio = ptr(*r.PayoutRatio * 100)
		}
	}

	if current.FreeCashFlow != nil {
		r.FreeCashFlow = ptr(float64(*current.FreeCashFlow))
	}

	if current.CAPEX != nil {
		r.CAPEX = ptr(float64(*current.CAPEX))
	}

	if current.EBITDA != nil {
		r.EBITDA = ptr(float64(*current.EBITDA))
	}

	if current.NetDebt != nil {
		r.NetDebt = ptr(float64(*current.NetDebt))
	}

	if current.WorkingCapital != nil {
		r.WorkingCapital = ptr(float64(*current.WorkingCapital))
	}

	if previous != nil {
		if current.Revenue != nil && previous.Revenue != nil {
			r.RevenueGrowth = pctGrowth(float64(*current.Revenue), float64(*previous.Revenue))
		}

		if current.NetProfit != nil && previous.NetProfit != nil {
			r.EarningsGrowth = pctGrowth(float64(*current.NetProfit), float64(*previous.NetProfit))
		}

		if current.EBITDA != nil && previous.EBITDA != nil {
			r.EBITDAGrowth = pctGrowth(float64(*current.EBITDA), float64(*previous.EBITDA))
		}

		if current.FreeCashFlow != nil && previous.FreeCashFlow != nil {
			r.FCFGrowth = pctGrowth(float64(*current.FreeCashFlow), float64(*previous.FreeCashFlow))
		}
	}

	if r.PriceToEarnings != nil && r.EarningsGrowth != nil && *r.EarningsGrowth > 0 {
		r.PEG = safeDiv(*r.PriceToEarnings, *r.EarningsGrowth)
	}

	return r
}
