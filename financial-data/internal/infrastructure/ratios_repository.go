package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RatiosRepository struct {
	pool *pgxpool.Pool
}

func NewRatiosRepository(pool *pgxpool.Pool) *RatiosRepository {
	return &RatiosRepository{pool: pool}
}

func (r *RatiosRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Ratios, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}

	query := `
		SELECT
			price_to_earnings, price_to_book, price_to_cash_flow, ev_to_ebitda, ev_to_sales, ev_to_fcf, peg,
			roe, roa, roic, gross_profit_margin, operating_profit_margin, net_profit_margin,
			current_ratio, quick_ratio,
			net_debt_to_ebitda, debt_to_equity, interest_coverage_ratio,
			income_quality, asset_turnover, inventory_turnover, receivables_turnover,
			eps, book_value_per_share, cash_flow_per_share, dividend_per_share, dividend_yield, payout_ratio,
			enterprise_value, market_cap, free_cash_flow, capex, ebitda, net_debt, working_capital,
			revenue_growth, earnings_growth, ebitda_growth, fcf_growth
		FROM ratios
		WHERE ticker = $1
	`

	ratios := &domain.Ratios{}
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&ratios.PriceToEarnings, &ratios.PriceToBook, &ratios.PriceToCashFlow, &ratios.EVToEBITDA, &ratios.EVToSales, &ratios.EVToFCF, &ratios.PEG,
		&ratios.ROE, &ratios.ROA, &ratios.ROIC, &ratios.GrossProfitMargin, &ratios.OperatingProfitMargin, &ratios.NetProfitMargin,
		&ratios.CurrentRatio, &ratios.QuickRatio,
		&ratios.NetDebtToEBITDA, &ratios.DebtToEquity, &ratios.InterestCoverageRatio,
		&ratios.IncomeQuality, &ratios.AssetTurnover, &ratios.InventoryTurnover, &ratios.ReceivablesTurnover,
		&ratios.EPS, &ratios.BookValuePerShare, &ratios.CashFlowPerShare, &ratios.DividendPerShare, &ratios.DividendYield, &ratios.PayoutRatio,
		&ratios.EnterpriseValue, &ratios.MarketCap, &ratios.FreeCashFlow, &ratios.CAPEX, &ratios.EBITDA, &ratios.NetDebt, &ratios.WorkingCapital,
		&ratios.RevenueGrowth, &ratios.EarningsGrowth, &ratios.EBITDAGrowth, &ratios.FCFGrowth,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ratios not found for ticker %s", ticker)
		}
		return nil, fmt.Errorf("failed to get ratios: %w", err)
	}

	return ratios, nil
}

func (r *RatiosRepository) GetBySector(ctx context.Context, sector domain.Sector) (*domain.Ratios, error) {
	if !sector.IsValid() {
		return nil, fmt.Errorf("invalid sector: %d", sector)
	}

	query := `
		SELECT
			AVG(price_to_earnings), AVG(price_to_book), AVG(price_to_cash_flow), AVG(ev_to_ebitda), AVG(ev_to_sales), AVG(ev_to_fcf), AVG(peg),
			AVG(roe), AVG(roa), AVG(roic), AVG(gross_profit_margin), AVG(operating_profit_margin), AVG(net_profit_margin),
			AVG(current_ratio), AVG(quick_ratio),
			AVG(net_debt_to_ebitda), AVG(debt_to_equity), AVG(interest_coverage_ratio),
			AVG(income_quality), AVG(asset_turnover), AVG(inventory_turnover), AVG(receivables_turnover),
			AVG(eps), AVG(book_value_per_share), AVG(cash_flow_per_share), AVG(dividend_per_share), AVG(dividend_yield), AVG(payout_ratio),
			AVG(enterprise_value), AVG(market_cap), AVG(free_cash_flow), AVG(capex), AVG(ebitda), AVG(net_debt), AVG(working_capital),
			AVG(revenue_growth), AVG(earnings_growth), AVG(ebitda_growth), AVG(fcf_growth)
		FROM ratios
		WHERE sector = $1
	`

	ratios := &domain.Ratios{}
	err := r.pool.QueryRow(ctx, query, sector).Scan(
		&ratios.PriceToEarnings, &ratios.PriceToBook, &ratios.PriceToCashFlow, &ratios.EVToEBITDA, &ratios.EVToSales, &ratios.EVToFCF, &ratios.PEG,
		&ratios.ROE, &ratios.ROA, &ratios.ROIC, &ratios.GrossProfitMargin, &ratios.OperatingProfitMargin, &ratios.NetProfitMargin,
		&ratios.CurrentRatio, &ratios.QuickRatio,
		&ratios.NetDebtToEBITDA, &ratios.DebtToEquity, &ratios.InterestCoverageRatio,
		&ratios.IncomeQuality, &ratios.AssetTurnover, &ratios.InventoryTurnover, &ratios.ReceivablesTurnover,
		&ratios.EPS, &ratios.BookValuePerShare, &ratios.CashFlowPerShare, &ratios.DividendPerShare, &ratios.DividendYield, &ratios.PayoutRatio,
		&ratios.EnterpriseValue, &ratios.MarketCap, &ratios.FreeCashFlow, &ratios.CAPEX, &ratios.EBITDA, &ratios.NetDebt, &ratios.WorkingCapital,
		&ratios.RevenueGrowth, &ratios.EarningsGrowth, &ratios.EBITDAGrowth, &ratios.FCFGrowth,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get average ratios for sector: %w", err)
	}

	return ratios, nil
}

func (r *RatiosRepository) Create(ctx context.Context, ticker string, sector domain.Sector, ratios *domain.Ratios) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if !sector.IsValid() {
		return fmt.Errorf("invalid sector: %d", sector)
	}
	if ratios == nil {
		return fmt.Errorf("ratios is nil")
	}

	query := `
		INSERT INTO ratios (
			ticker, sector,
			price_to_earnings, price_to_book, price_to_cash_flow, ev_to_ebitda, ev_to_sales, ev_to_fcf, peg,
			roe, roa, roic, gross_profit_margin, operating_profit_margin, net_profit_margin,
			current_ratio, quick_ratio,
			net_debt_to_ebitda, debt_to_equity, interest_coverage_ratio,
			income_quality, asset_turnover, inventory_turnover, receivables_turnover,
			eps, book_value_per_share, cash_flow_per_share, dividend_per_share, dividend_yield, payout_ratio,
			enterprise_value, market_cap, free_cash_flow, capex, ebitda, net_debt, working_capital,
			revenue_growth, earnings_growth, ebitda_growth, fcf_growth
		) VALUES (
			$1, $2,
			$3, $4, $5, $6, $7, $8, $9,
			$10, $11, $12, $13, $14, $15,
			$16, $17,
			$18, $19, $20,
			$21, $22, $23, $24,
			$25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35, $36, $37,
			$38, $39, $40, $41
		)
	`

	_, err := r.pool.Exec(ctx, query,
		ticker, sector,
		ratios.PriceToEarnings, ratios.PriceToBook, ratios.PriceToCashFlow, ratios.EVToEBITDA, ratios.EVToSales, ratios.EVToFCF, ratios.PEG,
		ratios.ROE, ratios.ROA, ratios.ROIC, ratios.GrossProfitMargin, ratios.OperatingProfitMargin, ratios.NetProfitMargin,
		ratios.CurrentRatio, ratios.QuickRatio,
		ratios.NetDebtToEBITDA, ratios.DebtToEquity, ratios.InterestCoverageRatio,
		ratios.IncomeQuality, ratios.AssetTurnover, ratios.InventoryTurnover, ratios.ReceivablesTurnover,
		ratios.EPS, ratios.BookValuePerShare, ratios.CashFlowPerShare, ratios.DividendPerShare, ratios.DividendYield, ratios.PayoutRatio,
		ratios.EnterpriseValue, ratios.MarketCap, ratios.FreeCashFlow, ratios.CAPEX, ratios.EBITDA, ratios.NetDebt, ratios.WorkingCapital,
		ratios.RevenueGrowth, ratios.EarningsGrowth, ratios.EBITDAGrowth, ratios.FCFGrowth,
	)

	if err != nil {
		return fmt.Errorf("failed to create ratios: %w", err)
	}

	return nil
}

func (r *RatiosRepository) Update(ctx context.Context, ticker string, ratios *domain.Ratios) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if ratios == nil {
		return fmt.Errorf("ratios is nil")
	}

	query := `
		UPDATE ratios SET
			price_to_earnings = $2, price_to_book = $3, price_to_cash_flow = $4, ev_to_ebitda = $5, ev_to_sales = $6, ev_to_fcf = $7, peg = $8,
			roe = $9, roa = $10, roic = $11, gross_profit_margin = $12, operating_profit_margin = $13, net_profit_margin = $14,
			current_ratio = $15, quick_ratio = $16,
			net_debt_to_ebitda = $17, debt_to_equity = $18, interest_coverage_ratio = $19,
			income_quality = $20, asset_turnover = $21, inventory_turnover = $22, receivables_turnover = $23,
			eps = $24, book_value_per_share = $25, cash_flow_per_share = $26, dividend_per_share = $27, dividend_yield = $28, payout_ratio = $29,
			enterprise_value = $30, market_cap = $31, free_cash_flow = $32, capex = $33, ebitda = $34, net_debt = $35, working_capital = $36,
			revenue_growth = $37, earnings_growth = $38, ebitda_growth = $39, fcf_growth = $40
		WHERE ticker = $1
	`

	result, err := r.pool.Exec(ctx, query,
		ticker,
		ratios.PriceToEarnings, ratios.PriceToBook, ratios.PriceToCashFlow, ratios.EVToEBITDA, ratios.EVToSales, ratios.EVToFCF, ratios.PEG,
		ratios.ROE, ratios.ROA, ratios.ROIC, ratios.GrossProfitMargin, ratios.OperatingProfitMargin, ratios.NetProfitMargin,
		ratios.CurrentRatio, ratios.QuickRatio,
		ratios.NetDebtToEBITDA, ratios.DebtToEquity, ratios.InterestCoverageRatio,
		ratios.IncomeQuality, ratios.AssetTurnover, ratios.InventoryTurnover, ratios.ReceivablesTurnover,
		ratios.EPS, ratios.BookValuePerShare, ratios.CashFlowPerShare, ratios.DividendPerShare, ratios.DividendYield, ratios.PayoutRatio,
		ratios.EnterpriseValue, ratios.MarketCap, ratios.FreeCashFlow, ratios.CAPEX, ratios.EBITDA, ratios.NetDebt, ratios.WorkingCapital,
		ratios.RevenueGrowth, ratios.EarningsGrowth, ratios.EBITDAGrowth, ratios.FCFGrowth,
	)

	if err != nil {
		return fmt.Errorf("failed to update ratios: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ratios not found for ticker %s", ticker)
	}

	return nil
}

func (r *RatiosRepository) Delete(ctx context.Context, ticker string) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}

	query := `DELETE FROM ratios WHERE ticker = $1`

	result, err := r.pool.Exec(ctx, query, ticker)
	if err != nil {
		return fmt.Errorf("failed to delete ratios: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ratios not found for ticker %s", ticker)
	}

	return nil
}
