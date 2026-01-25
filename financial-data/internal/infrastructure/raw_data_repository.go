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

func (r *RawDataRepository) GetByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}
	if year < 1900 || year > 2100 {
		return nil, fmt.Errorf("invalid year: %d", year)
	}
	if !period.IsValid() {
		return nil, fmt.Errorf("invalid period: %s", period)
	}

	query := `
		SELECT
			ticker, year, period,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			ebit, ebitda, interest_expense, tax_expense, net_profit,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			equity, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
			shares_outstanding, market_cap,
			working_capital, capital_employed, enterprise_value, net_debt
		FROM metrics
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	rawData := &domain.RawData{}
	err := r.pool.QueryRow(ctx, query, ticker, year, period).Scan(
		&rawData.Ticker, &rawData.Year, &rawData.Period,
		&rawData.Revenue, &rawData.CostOfRevenue, &rawData.GrossProfit, &rawData.OperatingExpenses,
		&rawData.EBIT, &rawData.EBITDA, &rawData.InterestExpense, &rawData.TaxExpense, &rawData.NetProfit,
		&rawData.TotalAssets, &rawData.CurrentAssets, &rawData.CashAndEquivalents, &rawData.Inventories, &rawData.Receivables,
		&rawData.TotalLiabilities, &rawData.CurrentLiabilities, &rawData.Debt, &rawData.LongTermDebt, &rawData.ShortTermDebt,
		&rawData.Equity, &rawData.RetainedEarnings,
		&rawData.OperatingCashFlow, &rawData.InvestingCashFlow, &rawData.FinancingCashFlow, &rawData.CAPEX, &rawData.FreeCashFlow,
		&rawData.SharesOutstanding, &rawData.MarketCap,
		&rawData.WorkingCapital, &rawData.CapitalEmployed, &rawData.EnterpriseValue, &rawData.NetDebt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("metrics not found for ticker %s, year %d, period %s", ticker, year, period)
		}
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return rawData, nil
}

func (r *RawDataRepository) GetLatestByTicker(ctx context.Context, ticker string) (*domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}

	query := `
		SELECT
			ticker, year, period,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			ebit, ebitda, interest_expense, tax_expense, net_profit,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			equity, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
			shares_outstanding, market_cap,
			working_capital, capital_employed, enterprise_value, net_debt
		FROM metrics
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
		LIMIT 1
	`

	rawData := &domain.RawData{}
	err := r.pool.QueryRow(ctx, query, ticker).Scan(
		&rawData.Ticker, &rawData.Year, &rawData.Period,
		&rawData.Revenue, &rawData.CostOfRevenue, &rawData.GrossProfit, &rawData.OperatingExpenses,
		&rawData.EBIT, &rawData.EBITDA, &rawData.InterestExpense, &rawData.TaxExpense, &rawData.NetProfit,
		&rawData.TotalAssets, &rawData.CurrentAssets, &rawData.CashAndEquivalents, &rawData.Inventories, &rawData.Receivables,
		&rawData.TotalLiabilities, &rawData.CurrentLiabilities, &rawData.Debt, &rawData.LongTermDebt, &rawData.ShortTermDebt,
		&rawData.Equity, &rawData.RetainedEarnings,
		&rawData.OperatingCashFlow, &rawData.InvestingCashFlow, &rawData.FinancingCashFlow, &rawData.CAPEX, &rawData.FreeCashFlow,
		&rawData.SharesOutstanding, &rawData.MarketCap,
		&rawData.WorkingCapital, &rawData.CapitalEmployed, &rawData.EnterpriseValue, &rawData.NetDebt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no metrics found for ticker %s", ticker)
		}
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}

	return rawData, nil
}

func (r *RawDataRepository) GetHistoryByTicker(ctx context.Context, ticker string) ([]domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty")
	}

	query := `
		SELECT
			ticker, year, period,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			ebit, ebitda, interest_expense, tax_expense, net_profit,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			equity, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
			shares_outstanding, market_cap,
			working_capital, capital_employed, enterprise_value, net_debt
		FROM metrics
		WHERE ticker = $1
		ORDER BY year DESC, period DESC
	`

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var history []domain.RawData
	for rows.Next() {
		var rawData domain.RawData
		err := rows.Scan(
			&rawData.Ticker, &rawData.Year, &rawData.Period,
			&rawData.Revenue, &rawData.CostOfRevenue, &rawData.GrossProfit, &rawData.OperatingExpenses,
			&rawData.EBIT, &rawData.EBITDA, &rawData.InterestExpense, &rawData.TaxExpense, &rawData.NetProfit,
			&rawData.TotalAssets, &rawData.CurrentAssets, &rawData.CashAndEquivalents, &rawData.Inventories, &rawData.Receivables,
			&rawData.TotalLiabilities, &rawData.CurrentLiabilities, &rawData.Debt, &rawData.LongTermDebt, &rawData.ShortTermDebt,
			&rawData.Equity, &rawData.RetainedEarnings,
			&rawData.OperatingCashFlow, &rawData.InvestingCashFlow, &rawData.FinancingCashFlow, &rawData.CAPEX, &rawData.FreeCashFlow,
			&rawData.SharesOutstanding, &rawData.MarketCap,
			&rawData.WorkingCapital, &rawData.CapitalEmployed, &rawData.EnterpriseValue, &rawData.NetDebt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metrics: %w", err)
		}
		history = append(history, rawData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating metrics: %w", err)
	}

	return history, nil
}

func (r *RawDataRepository) Create(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return fmt.Errorf("rawData is nil")
	}
	if rawData.Ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return fmt.Errorf("invalid year: %d", rawData.Year)
	}
	if !rawData.Period.IsValid() {
		return fmt.Errorf("invalid period: %s", rawData.Period)
	}

	query := `
		INSERT INTO metrics (
			ticker, year, period,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			ebit, ebitda, interest_expense, tax_expense, net_profit,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			equity, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow, capex, free_cash_flow,
			shares_outstanding, market_cap,
			working_capital, capital_employed, enterprise_value, net_debt
		) VALUES (
			$1, $2, $3,
			$4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22,
			$23, $24,
			$25, $26, $27, $28, $29,
			$30, $31,
			$32, $33, $34, $35
		)
	`

	_, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period,
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
		return fmt.Errorf("failed to create metrics: %w", err)
	}

	return nil
}

func (r *RawDataRepository) Update(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return fmt.Errorf("rawData is nil")
	}
	if rawData.Ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return fmt.Errorf("invalid year: %d", rawData.Year)
	}
	if !rawData.Period.IsValid() {
		return fmt.Errorf("invalid period: %s", rawData.Period)
	}

	query := `
		UPDATE metrics SET
			revenue = $4, cost_of_revenue = $5, gross_profit = $6, operating_expenses = $7,
			ebit = $8, ebitda = $9, interest_expense = $10, tax_expense = $11, net_profit = $12,
			total_assets = $13, current_assets = $14, cash_and_equivalents = $15, inventories = $16, receivables = $17,
			total_liabilities = $18, current_liabilities = $19, debt = $20, long_term_debt = $21, short_term_debt = $22,
			equity = $23, retained_earnings = $24,
			operating_cash_flow = $25, investing_cash_flow = $26, financing_cash_flow = $27, capex = $28, free_cash_flow = $29,
			shares_outstanding = $30, market_cap = $31,
			working_capital = $32, capital_employed = $33, enterprise_value = $34, net_debt = $35
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	result, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period,
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
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("metrics not found for ticker %s, year %d, period %s", rawData.Ticker, rawData.Year, rawData.Period)
	}

	return nil
}

func (r *RawDataRepository) Delete(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty")
	}
	if year < 1900 || year > 2100 {
		return fmt.Errorf("invalid year: %d", year)
	}
	if !period.IsValid() {
		return fmt.Errorf("invalid period: %s", period)
	}

	query := `DELETE FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return fmt.Errorf("failed to delete metrics: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("metrics not found for ticker %s, year %d, period %s", ticker, year, period)
	}

	return nil
}
