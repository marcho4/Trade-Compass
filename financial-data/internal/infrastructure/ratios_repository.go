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

const ratiosSelectColumns = `
	ticker, year, period,
	price_to_earnings, price_to_book, price_to_cash_flow, ev_to_ebitda, ev_to_sales, ev_to_fcf, peg,
	roe, roa, roic, gross_profit_margin, operating_profit_margin, net_profit_margin,
	current_ratio, quick_ratio,
	net_debt_to_ebitda, debt_to_equity, interest_coverage_ratio,
	income_quality, asset_turnover, inventory_turnover, receivables_turnover,
	eps, book_value_per_share, cash_flow_per_share, dividend_per_share, dividend_yield, payout_ratio,
	enterprise_value, market_cap, free_cash_flow, capex, ebitda, net_debt, working_capital,
	revenue_growth, earnings_growth, ebitda_growth, fcf_growth
`

func ratiosScanTargets(r *domain.Ratios) []any {
	return []any{
		&r.Ticker, &r.Year, &r.Period,
		&r.PriceToEarnings, &r.PriceToBook, &r.PriceToCashFlow, &r.EVToEBITDA, &r.EVToSales, &r.EVToFCF, &r.PEG,
		&r.ROE, &r.ROA, &r.ROIC, &r.GrossProfitMargin, &r.OperatingProfitMargin, &r.NetProfitMargin,
		&r.CurrentRatio, &r.QuickRatio,
		&r.NetDebtToEBITDA, &r.DebtToEquity, &r.InterestCoverageRatio,
		&r.IncomeQuality, &r.AssetTurnover, &r.InventoryTurnover, &r.ReceivablesTurnover,
		&r.EPS, &r.BookValuePerShare, &r.CashFlowPerShare, &r.DividendPerShare, &r.DividendYield, &r.PayoutRatio,
		&r.EnterpriseValue, &r.MarketCap, &r.FreeCashFlow, &r.CAPEX, &r.EBITDA, &r.NetDebt, &r.WorkingCapital,
		&r.RevenueGrowth, &r.EarningsGrowth, &r.EBITDAGrowth, &r.FCFGrowth,
	}
}

const ratiosValueColumns = `
	price_to_earnings, price_to_book, price_to_cash_flow, ev_to_ebitda, ev_to_sales, ev_to_fcf, peg,
	roe, roa, roic, gross_profit_margin, operating_profit_margin, net_profit_margin,
	current_ratio, quick_ratio,
	net_debt_to_ebitda, debt_to_equity, interest_coverage_ratio,
	income_quality, asset_turnover, inventory_turnover, receivables_turnover,
	eps, book_value_per_share, cash_flow_per_share, dividend_per_share, dividend_yield, payout_ratio,
	enterprise_value, market_cap, free_cash_flow, capex, ebitda, net_debt, working_capital,
	revenue_growth, earnings_growth, ebitda_growth, fcf_growth
`

func ratiosValueArgs(r *domain.Ratios) []any {
	return []any{
		r.PriceToEarnings, r.PriceToBook, r.PriceToCashFlow, r.EVToEBITDA, r.EVToSales, r.EVToFCF, r.PEG,
		r.ROE, r.ROA, r.ROIC, r.GrossProfitMargin, r.OperatingProfitMargin, r.NetProfitMargin,
		r.CurrentRatio, r.QuickRatio,
		r.NetDebtToEBITDA, r.DebtToEquity, r.InterestCoverageRatio,
		r.IncomeQuality, r.AssetTurnover, r.InventoryTurnover, r.ReceivablesTurnover,
		r.EPS, r.BookValuePerShare, r.CashFlowPerShare, r.DividendPerShare, r.DividendYield, r.PayoutRatio,
		r.EnterpriseValue, r.MarketCap, r.FreeCashFlow, r.CAPEX, r.EBITDA, r.NetDebt, r.WorkingCapital,
		r.RevenueGrowth, r.EarningsGrowth, r.EBITDAGrowth, r.FCFGrowth,
	}
}

func (r *RatiosRepository) GetByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.Ratios, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM ratios WHERE ticker = $1 AND year = $2 AND period = $3`, ratiosSelectColumns)

	ratios := &domain.Ratios{}
	err := r.pool.QueryRow(ctx, query, ticker, year, period).Scan(ratiosScanTargets(ratios)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ratios not found for ticker %s year %d period %s: %w", ticker, year, period, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get ratios: %w", err)
	}

	return ratios, nil
}

func (r *RatiosRepository) GetLatestByTicker(ctx context.Context, ticker string) (*domain.Ratios, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM ratios WHERE ticker = $1 ORDER BY year DESC, period DESC LIMIT 1`, ratiosSelectColumns)

	ratios := &domain.Ratios{}
	err := r.pool.QueryRow(ctx, query, ticker).Scan(ratiosScanTargets(ratios)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ratios not found for ticker %s: %w", ticker, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get latest ratios: %w", err)
	}

	return ratios, nil
}

func (r *RatiosRepository) GetHistoryByTicker(ctx context.Context, ticker string) ([]domain.Ratios, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM ratios WHERE ticker = $1 ORDER BY year DESC, period DESC`, ratiosSelectColumns)

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get ratios history: %w", err)
	}
	defer rows.Close()

	var result []domain.Ratios
	for rows.Next() {
		var ratios domain.Ratios
		if err := rows.Scan(ratiosScanTargets(&ratios)...); err != nil {
			return nil, fmt.Errorf("failed to scan ratios row: %w", err)
		}
		result = append(result, ratios)
	}

	return result, nil
}

func (r *RatiosRepository) GetBySector(ctx context.Context, sector domain.Sector) (*domain.Ratios, error) {
	if !sector.IsValid() {
		return nil, fmt.Errorf("invalid sector: %d: %w", sector, domain.ErrInvalidInput)
	}

	query := `
		SELECT
			'' AS ticker, 0 AS year, '' AS period,
			AVG(price_to_earnings), AVG(price_to_book), AVG(price_to_cash_flow), AVG(ev_to_ebitda), AVG(ev_to_sales), AVG(ev_to_fcf), AVG(peg),
			AVG(roe), AVG(roa), AVG(roic), AVG(gross_profit_margin), AVG(operating_profit_margin), AVG(net_profit_margin),
			AVG(current_ratio), AVG(quick_ratio),
			AVG(net_debt_to_ebitda), AVG(debt_to_equity), AVG(interest_coverage_ratio),
			AVG(income_quality), AVG(asset_turnover), AVG(inventory_turnover), AVG(receivables_turnover),
			AVG(eps), AVG(book_value_per_share), AVG(cash_flow_per_share), AVG(dividend_per_share), AVG(dividend_yield), AVG(payout_ratio),
			AVG(enterprise_value), AVG(market_cap), AVG(free_cash_flow), AVG(capex), AVG(ebitda), AVG(net_debt), AVG(working_capital),
			AVG(revenue_growth), AVG(earnings_growth), AVG(ebitda_growth), AVG(fcf_growth)
		FROM (
			SELECT DISTINCT ON (ticker) *
			FROM ratios
			WHERE sector = $1
			ORDER BY ticker, year DESC, period DESC
		) AS latest
	`

	ratios := &domain.Ratios{}
	err := r.pool.QueryRow(ctx, query, sector).Scan(ratiosScanTargets(ratios)...)
	if err != nil {
		return nil, fmt.Errorf("failed to get average ratios for sector: %w", err)
	}

	return ratios, nil
}

func (r *RatiosRepository) Create(ctx context.Context, sector domain.Sector, ratios *domain.Ratios) error {
	if ratios == nil {
		return fmt.Errorf("ratios is nil: %w", domain.ErrInvalidInput)
	}
	if ratios.Ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if !sector.IsValid() {
		return fmt.Errorf("invalid sector: %d: %w", sector, domain.ErrInvalidInput)
	}

	query := `
		INSERT INTO ratios (
			ticker, year, period, sector,
			` + ratiosValueColumns + `
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9, $10, $11,
			$12, $13, $14, $15, $16, $17,
			$18, $19,
			$20, $21, $22,
			$23, $24, $25, $26,
			$27, $28, $29, $30, $31, $32,
			$33, $34, $35, $36, $37, $38, $39,
			$40, $41, $42, $43
		)
	`

	args := []any{ratios.Ticker, ratios.Year, ratios.Period, sector}
	args = append(args, ratiosValueArgs(ratios)...)

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create ratios: %w", err)
	}

	return nil
}

func (r *RatiosRepository) Update(ctx context.Context, ratios *domain.Ratios) error {
	if ratios == nil {
		return fmt.Errorf("ratios is nil: %w", domain.ErrInvalidInput)
	}
	if ratios.Ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := `
		UPDATE ratios SET
			price_to_earnings = $4, price_to_book = $5, price_to_cash_flow = $6, ev_to_ebitda = $7, ev_to_sales = $8, ev_to_fcf = $9, peg = $10,
			roe = $11, roa = $12, roic = $13, gross_profit_margin = $14, operating_profit_margin = $15, net_profit_margin = $16,
			current_ratio = $17, quick_ratio = $18,
			net_debt_to_ebitda = $19, debt_to_equity = $20, interest_coverage_ratio = $21,
			income_quality = $22, asset_turnover = $23, inventory_turnover = $24, receivables_turnover = $25,
			eps = $26, book_value_per_share = $27, cash_flow_per_share = $28, dividend_per_share = $29, dividend_yield = $30, payout_ratio = $31,
			enterprise_value = $32, market_cap = $33, free_cash_flow = $34, capex = $35, ebitda = $36, net_debt = $37, working_capital = $38,
			revenue_growth = $39, earnings_growth = $40, ebitda_growth = $41, fcf_growth = $42
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	args := []any{ratios.Ticker, ratios.Year, ratios.Period}
	args = append(args, ratiosValueArgs(ratios)...)

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update ratios: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ratios not found for ticker %s year %d period %s: %w", ratios.Ticker, ratios.Year, ratios.Period, domain.ErrNotFound)
	}

	return nil
}

func (r *RatiosRepository) Delete(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := `DELETE FROM ratios WHERE ticker = $1 AND year = $2 AND period = $3`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return fmt.Errorf("failed to delete ratios: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("ratios not found for ticker %s: %w", ticker, domain.ErrNotFound)
	}

	return nil
}
