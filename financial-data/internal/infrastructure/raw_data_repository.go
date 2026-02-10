package infrastructure

import (
	"context"
	"errors"
	"financial_data/internal/domain"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RawDataRepository struct {
	pool *pgxpool.Pool
}

func NewRawDataRepository(pool *pgxpool.Pool) *RawDataRepository {
	return &RawDataRepository{pool: pool}
}

const rawDataSelectColumns = `
	ticker, year, period, status,
	revenue, cost_of_revenue, gross_profit, operating_expenses,
	ebit, ebitda, interest_expense, tax_expense, net_profit,
	total_assets, current_assets, cash_and_equivalents, inventories, receivables,
	total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
	equity, retained_earnings,
	operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
	shares_outstanding, market_cap,
	working_capital, capital_employed, enterprise_value, net_debt
`

func scanRawData(row pgx.Row) (*domain.RawData, error) {
	rd := &domain.RawData{}
	err := row.Scan(
		&rd.Ticker, &rd.Year, &rd.Period, &rd.Status,
		&rd.Revenue, &rd.CostOfRevenue, &rd.GrossProfit, &rd.OperatingExpenses,
		&rd.EBIT, &rd.EBITDA, &rd.InterestExpense, &rd.TaxExpense, &rd.NetProfit,
		&rd.TotalAssets, &rd.CurrentAssets, &rd.CashAndEquivalents, &rd.Inventories, &rd.Receivables,
		&rd.TotalLiabilities, &rd.CurrentLiabilities, &rd.Debt, &rd.LongTermDebt, &rd.ShortTermDebt,
		&rd.Equity, &rd.RetainedEarnings,
		&rd.OperatingCashFlow, &rd.InvestingCashFlow, &rd.FinancingCashFlow, &rd.CAPEX, &rd.FreeCashFlow,
		&rd.SharesOutstanding, &rd.MarketCap,
		&rd.WorkingCapital, &rd.CapitalEmployed, &rd.EnterpriseValue, &rd.NetDebt,
	)
	return rd, err
}

func scanRawDataRows(rows pgx.Rows) ([]domain.RawData, error) {
	var result []domain.RawData
	for rows.Next() {
		var rd domain.RawData
		err := rows.Scan(
			&rd.Ticker, &rd.Year, &rd.Period, &rd.Status,
			&rd.Revenue, &rd.CostOfRevenue, &rd.GrossProfit, &rd.OperatingExpenses,
			&rd.EBIT, &rd.EBITDA, &rd.InterestExpense, &rd.TaxExpense, &rd.NetProfit,
			&rd.TotalAssets, &rd.CurrentAssets, &rd.CashAndEquivalents, &rd.Inventories, &rd.Receivables,
			&rd.TotalLiabilities, &rd.CurrentLiabilities, &rd.Debt, &rd.LongTermDebt, &rd.ShortTermDebt,
			&rd.Equity, &rd.RetainedEarnings,
			&rd.OperatingCashFlow, &rd.InvestingCashFlow, &rd.FinancingCashFlow, &rd.CAPEX, &rd.FreeCashFlow,
			&rd.SharesOutstanding, &rd.MarketCap,
			&rd.WorkingCapital, &rd.CapitalEmployed, &rd.EnterpriseValue, &rd.NetDebt,
		)
		if err != nil {
			return nil, NewDbError(fmt.Sprintf("failed to scan metrics: %v", err), 0)
		}
		result = append(result, rd)
	}
	if err := rows.Err(); err != nil {
		return nil, NewDbError(fmt.Sprintf("error iterating metrics: %v", err), 0)
	}
	return result, nil
}

func (r *RawDataRepository) GetByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}
	if year < 1900 || year > 2100 {
		return nil, NewDbError(fmt.Sprintf("invalid year: %d", year), 0)
	}
	if !period.IsValid() {
		return nil, NewDbError(fmt.Sprintf("invalid period: %s", period), 0)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'confirmed'`, rawDataSelectColumns)

	rd, err := scanRawData(r.pool.QueryRow(ctx, query, ticker, year, period))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError(fmt.Sprintf("metrics not found for ticker %s, year %d, period %s", ticker, year, period), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get metrics: %v", err), 0)
	}
	return rd, nil
}

func (r *RawDataRepository) GetLatestByTicker(ctx context.Context, ticker string) (*domain.RawData, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'confirmed' ORDER BY year DESC, period DESC LIMIT 1`, rawDataSelectColumns)

	rd, err := scanRawData(r.pool.QueryRow(ctx, query, ticker))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NewDbError(fmt.Sprintf("no metrics found for ticker %s", ticker), 0)
		}
		return nil, NewDbError(fmt.Sprintf("failed to get latest metrics: %v", err), 0)
	}
	return rd, nil
}

func (r *RawDataRepository) GetHistoryByTicker(ctx context.Context, ticker string) ([]domain.RawData, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'confirmed' ORDER BY year DESC, period DESC`, rawDataSelectColumns)

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query metrics: %v", err), 0)
	}
	defer rows.Close()

	return scanRawDataRows(rows)
}

func (r *RawDataRepository) GetDraftByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'draft'`, rawDataSelectColumns)

	rd, err := scanRawData(r.pool.QueryRow(ctx, query, ticker, year, period))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, NewDbError(fmt.Sprintf("failed to get draft: %v", err), 0)
	}
	return rd, nil
}

func (r *RawDataRepository) GetDraftsByTicker(ctx context.Context, ticker string) ([]domain.RawData, error) {
	if ticker == "" {
		return nil, NewDbError("ticker is empty", 0)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'draft' ORDER BY year DESC, period DESC`, rawDataSelectColumns)

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, NewDbError(fmt.Sprintf("failed to query drafts: %v", err), 0)
	}
	defer rows.Close()

	return scanRawDataRows(rows)
}

func (r *RawDataRepository) ConfirmDraft(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return NewDbError("ticker is empty", 0)
	}

	query := `UPDATE metrics SET status = 'confirmed', updated_at = NOW() WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'draft'`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to confirm draft: %v", err), 0)
	}
	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("draft not found for ticker %s, year %d, period %s", ticker, year, period), 0)
	}
	return nil
}

func (r *RawDataRepository) Create(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return NewDbError("rawData is nil", 0)
	}
	if rawData.Ticker == "" {
		return NewDbError("ticker is empty", 0)
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return NewDbError(fmt.Sprintf("invalid year: %d", rawData.Year), 0)
	}
	if !rawData.Period.IsValid() {
		return NewDbError(fmt.Sprintf("invalid period: %s", rawData.Period), 0)
	}

	status := rawData.Status
	if !status.IsValid() {
		status = domain.StatusConfirmed
	}

	query := `
		INSERT INTO metrics (
			ticker, year, period, status,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			ebit, ebitda, interest_expense, tax_expense, net_profit,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			equity, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
			shares_outstanding, market_cap,
			working_capital, capital_employed, enterprise_value, net_debt
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18,
			$19, $20, $21, $22, $23,
			$24, $25,
			$26, $27, $28, $29, $30,
			$31, $32,
			$33, $34, $35, $36
		)
	`

	_, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period, status,
		rawData.Revenue, rawData.CostOfRevenue, rawData.GrossProfit, rawData.OperatingExpenses,
		rawData.EBIT, rawData.EBITDA, rawData.InterestExpense, rawData.TaxExpense, rawData.NetProfit,
		rawData.TotalAssets, rawData.CurrentAssets, rawData.CashAndEquivalents, rawData.Inventories, rawData.Receivables,
		rawData.TotalLiabilities, rawData.CurrentLiabilities, rawData.Debt, rawData.LongTermDebt, rawData.ShortTermDebt,
		rawData.Equity, rawData.RetainedEarnings,
		rawData.OperatingCashFlow, rawData.InvestingCashFlow, rawData.FinancingCashFlow, rawData.CAPEX, rawData.FreeCashFlow,
		rawData.SharesOutstanding, rawData.MarketCap,
		rawData.WorkingCapital, rawData.CapitalEmployed, rawData.EnterpriseValue, rawData.NetDebt,
	)

	if err != nil {
		return NewDbError(fmt.Sprintf("failed to create metrics: %v", err), 0)
	}

	return nil
}

func (r *RawDataRepository) Update(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return NewDbError("rawData is nil", 0)
	}
	if rawData.Ticker == "" {
		return NewDbError("ticker is empty", 0)
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return NewDbError(fmt.Sprintf("invalid year: %d", rawData.Year), 0)
	}
	if !rawData.Period.IsValid() {
		return NewDbError(fmt.Sprintf("invalid period: %s", rawData.Period), 0)
	}

	status := rawData.Status
	if !status.IsValid() {
		status = domain.StatusConfirmed
	}

	query := `
		UPDATE metrics SET
			status = $4,
			revenue = $5, cost_of_revenue = $6, gross_profit = $7, operating_expenses = $8,
			ebit = $9, ebitda = $10, interest_expense = $11, tax_expense = $12, net_profit = $13,
			total_assets = $14, current_assets = $15, cash_and_equivalents = $16, inventories = $17, receivables = $18,
			total_liabilities = $19, current_liabilities = $20, debt = $21, long_term_debt = $22, short_term_debt = $23,
			equity = $24, retained_earnings = $25,
			operating_cash_flow = $26, investing_cash_flow = $27, financing_cash_flow = $28, capex = $29, free_cash_flow = $30,
			shares_outstanding = $31, market_cap = $32,
			working_capital = $33, capital_employed = $34, enterprise_value = $35, net_debt = $36,
			updated_at = NOW()
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	result, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period, status,
		rawData.Revenue, rawData.CostOfRevenue, rawData.GrossProfit, rawData.OperatingExpenses,
		rawData.EBIT, rawData.EBITDA, rawData.InterestExpense, rawData.TaxExpense, rawData.NetProfit,
		rawData.TotalAssets, rawData.CurrentAssets, rawData.CashAndEquivalents, rawData.Inventories, rawData.Receivables,
		rawData.TotalLiabilities, rawData.CurrentLiabilities, rawData.Debt, rawData.LongTermDebt, rawData.ShortTermDebt,
		rawData.Equity, rawData.RetainedEarnings,
		rawData.OperatingCashFlow, rawData.InvestingCashFlow, rawData.FinancingCashFlow, rawData.CAPEX, rawData.FreeCashFlow,
		rawData.SharesOutstanding, rawData.MarketCap,
		rawData.WorkingCapital, rawData.CapitalEmployed, rawData.EnterpriseValue, rawData.NetDebt,
	)

	if err != nil {
		return NewDbError(fmt.Sprintf("failed to update metrics: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("metrics not found for ticker %s, year %d, period %s", rawData.Ticker, rawData.Year, rawData.Period), 0)
	}

	return nil
}

func (r *RawDataRepository) Delete(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return NewDbError("ticker is empty", 0)
	}
	if year < 1900 || year > 2100 {
		return NewDbError(fmt.Sprintf("invalid year: %d", year), 0)
	}
	if !period.IsValid() {
		return NewDbError(fmt.Sprintf("invalid period: %s", period), 0)
	}

	query := `DELETE FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return NewDbError(fmt.Sprintf("failed to delete metrics: %v", err), 0)
	}

	if result.RowsAffected() == 0 {
		return NewDbError(fmt.Sprintf("metrics not found for ticker %s, year %d, period %s", ticker, year, period), 0)
	}

	return nil
}
