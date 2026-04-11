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
	ticker, year, period, status, report_units,
	revenue, cost_of_revenue, gross_profit, operating_expenses,
	other_income, other_expenses,
	ebit, ebitda, depreciation,
	interest_income, interest_expense,
	profit_before_tax, tax_expense, net_profit, net_profit_parent, basic_eps,
	total_assets, current_assets, cash_and_equivalents, inventories, receivables,
	fixed_assets, right_of_use_assets, intangible_assets, goodwill, total_non_current_assets,
	total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
	lt_lease_liabilities, st_lease_liabilities, trade_payables,
	equity, equity_parent, treasury_shares, retained_earnings,
	operating_cash_flow, investing_cash_flow, financing_cash_flow,
	capex, free_cash_flow, dividends_paid, lease_payments,
	acquisitions_net, interest_paid, debt_proceeds, debt_repayments,
	shares_outstanding, market_cap, enterprise_value,
	working_capital, capital_employed, net_debt,
	interest_on_leases, interest_on_loans,
	company_type, net_interest_income, commission_income, commission_expense, net_commission_income, credit_loss_provision
`

func rawDataScanTargets(rd *domain.RawData) []any {
	return []any{
		&rd.Ticker, &rd.Year, &rd.Period, &rd.Status, &rd.ReportUnits,
		&rd.Revenue, &rd.CostOfRevenue, &rd.GrossProfit, &rd.OperatingExpenses,
		&rd.OtherIncome, &rd.OtherExpenses,
		&rd.EBIT, &rd.EBITDA, &rd.Depreciation,
		&rd.InterestIncome, &rd.InterestExpense,
		&rd.ProfitBeforeTax, &rd.TaxExpense, &rd.NetProfit, &rd.NetProfitParent, &rd.BasicEPS,
		&rd.TotalAssets, &rd.CurrentAssets, &rd.CashAndEquivalents, &rd.Inventories, &rd.Receivables,
		&rd.FixedAssets, &rd.RightOfUseAssets, &rd.IntangibleAssets, &rd.Goodwill, &rd.TotalNonCurrentAssets,
		&rd.TotalLiabilities, &rd.CurrentLiabilities, &rd.Debt, &rd.LongTermDebt, &rd.ShortTermDebt,
		&rd.LtLeaseLiabilities, &rd.StLeaseLiabilities, &rd.TradePayables,
		&rd.Equity, &rd.EquityParent, &rd.TreasuryShares, &rd.RetainedEarnings,
		&rd.OperatingCashFlow, &rd.InvestingCashFlow, &rd.FinancingCashFlow,
		&rd.CAPEX, &rd.FreeCashFlow, &rd.DividendsPaid, &rd.LeasePayments,
		&rd.AcquisitionsNet, &rd.InterestPaid, &rd.DebtProceeds, &rd.DebtRepayments,
		&rd.SharesOutstanding, &rd.MarketCap, &rd.EnterpriseValue,
		&rd.WorkingCapital, &rd.CapitalEmployed, &rd.NetDebt,
		&rd.InterestOnLeases, &rd.InterestOnLoans,
		&rd.CompanyType, &rd.NetInterestIncome, &rd.CommissionIncome, &rd.CommissionExpense, &rd.NetCommissionIncome, &rd.CreditLossProvision,
	}
}

func (r *RawDataRepository) GetByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if year < 1900 || year > 2100 {
		return nil, fmt.Errorf("invalid year: %d: %w", year, domain.ErrInvalidInput)
	}
	if !period.IsValid() {
		return nil, fmt.Errorf("invalid period: %s: %w", period, domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'confirmed'`, rawDataSelectColumns)

	rd := &domain.RawData{}
	err := r.pool.QueryRow(ctx, query, ticker, year, period).Scan(rawDataScanTargets(rd)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("metrics not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	return rd, nil
}

func (r *RawDataRepository) GetLatestByTicker(ctx context.Context, ticker string) (*domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'confirmed' ORDER BY year DESC, period DESC LIMIT 1`, rawDataSelectColumns)

	rd := &domain.RawData{}
	err := r.pool.QueryRow(ctx, query, ticker).Scan(rawDataScanTargets(rd)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no metrics found for ticker %s: %w", ticker, domain.ErrNotFound)
		}
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}
	return rd, nil
}

func (r *RawDataRepository) GetHistoryByTicker(ctx context.Context, ticker string) ([]domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'confirmed' ORDER BY year DESC, period DESC`, rawDataSelectColumns)

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var result []domain.RawData
	for rows.Next() {
		var rd domain.RawData
		if err := rows.Scan(rawDataScanTargets(&rd)...); err != nil {
			return nil, fmt.Errorf("failed to scan metrics: %w", err)
		}
		result = append(result, rd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating metrics: %w", err)
	}
	return result, nil
}

func (r *RawDataRepository) GetDraftByTickerAndPeriod(ctx context.Context, ticker string, year int, period domain.ReportPeriod) (*domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'draft'`, rawDataSelectColumns)

	rd := &domain.RawData{}
	err := r.pool.QueryRow(ctx, query, ticker, year, period).Scan(rawDataScanTargets(rd)...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get draft: %w", err)
	}
	return rd, nil
}

func (r *RawDataRepository) GetDraftsByTicker(ctx context.Context, ticker string) ([]domain.RawData, error) {
	if ticker == "" {
		return nil, fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := fmt.Sprintf(`SELECT %s FROM metrics WHERE ticker = $1 AND status = 'draft' ORDER BY year DESC, period DESC`, rawDataSelectColumns)

	rows, err := r.pool.Query(ctx, query, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to query drafts: %w", err)
	}
	defer rows.Close()

	var result []domain.RawData
	for rows.Next() {
		var rd domain.RawData
		if err := rows.Scan(rawDataScanTargets(&rd)...); err != nil {
			return nil, fmt.Errorf("failed to scan draft: %w", err)
		}
		result = append(result, rd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating drafts: %w", err)
	}
	return result, nil
}

func (r *RawDataRepository) ConfirmDraft(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}

	query := `UPDATE metrics SET status = 'confirmed', updated_at = NOW() WHERE ticker = $1 AND year = $2 AND period = $3 AND status = 'draft'`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return fmt.Errorf("failed to confirm draft: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("draft not found for ticker %s, year %d, period %s: %w", ticker, year, period, domain.ErrNotFound)
	}
	return nil
}

func (r *RawDataRepository) Create(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return fmt.Errorf("rawData is nil: %w", domain.ErrInvalidInput)
	}
	if rawData.Ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return fmt.Errorf("invalid year: %d: %w", rawData.Year, domain.ErrInvalidInput)
	}
	if !rawData.Period.IsValid() {
		return fmt.Errorf("invalid period: %s: %w", rawData.Period, domain.ErrInvalidInput)
	}

	status := rawData.Status
	if !status.IsValid() {
		status = domain.StatusConfirmed
	}

	query := `
		INSERT INTO metrics (
			ticker, year, period, status, report_units,
			revenue, cost_of_revenue, gross_profit, operating_expenses,
			other_income, other_expenses,
			ebit, ebitda, depreciation,
			interest_income, interest_expense,
			profit_before_tax, tax_expense, net_profit, net_profit_parent, basic_eps,
			total_assets, current_assets, cash_and_equivalents, inventories, receivables,
			fixed_assets, right_of_use_assets, intangible_assets, goodwill, total_non_current_assets,
			total_liabilities, current_liabilities, debt, long_term_debt, short_term_debt,
			lt_lease_liabilities, st_lease_liabilities, trade_payables,
			equity, equity_parent, treasury_shares, retained_earnings,
			operating_cash_flow, investing_cash_flow, financing_cash_flow,
			capex, free_cash_flow, dividends_paid, lease_payments,
			acquisitions_net, interest_paid, debt_proceeds, debt_repayments,
			shares_outstanding, market_cap, enterprise_value,
			working_capital, capital_employed, net_debt,
			interest_on_leases, interest_on_loans,
			company_type, net_interest_income, commission_income, commission_expense, net_commission_income, credit_loss_provision
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11,
			$12, $13, $14,
			$15, $16,
			$17, $18, $19, $20, $21,
			$22, $23, $24, $25, $26,
			$27, $28, $29, $30, $31,
			$32, $33, $34, $35, $36,
			$37, $38, $39,
			$40, $41, $42, $43,
			$44, $45, $46,
			$47, $48, $49, $50,
			$51, $52, $53, $54,
			$55, $56, $57,
			$58, $59, $60,
			$61, $62,
			$63, $64, $65, $66, $67, $68
		)
	`

	_, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period, status, rawData.ReportUnits,
		rawData.Revenue, rawData.CostOfRevenue, rawData.GrossProfit, rawData.OperatingExpenses,
		rawData.OtherIncome, rawData.OtherExpenses,
		rawData.EBIT, rawData.EBITDA, rawData.Depreciation,
		rawData.InterestIncome, rawData.InterestExpense,
		rawData.ProfitBeforeTax, rawData.TaxExpense, rawData.NetProfit, rawData.NetProfitParent, rawData.BasicEPS,
		rawData.TotalAssets, rawData.CurrentAssets, rawData.CashAndEquivalents, rawData.Inventories, rawData.Receivables,
		rawData.FixedAssets, rawData.RightOfUseAssets, rawData.IntangibleAssets, rawData.Goodwill, rawData.TotalNonCurrentAssets,
		rawData.TotalLiabilities, rawData.CurrentLiabilities, rawData.Debt, rawData.LongTermDebt, rawData.ShortTermDebt,
		rawData.LtLeaseLiabilities, rawData.StLeaseLiabilities, rawData.TradePayables,
		rawData.Equity, rawData.EquityParent, rawData.TreasuryShares, rawData.RetainedEarnings,
		rawData.OperatingCashFlow, rawData.InvestingCashFlow, rawData.FinancingCashFlow,
		rawData.CAPEX, rawData.FreeCashFlow, rawData.DividendsPaid, rawData.LeasePayments,
		rawData.AcquisitionsNet, rawData.InterestPaid, rawData.DebtProceeds, rawData.DebtRepayments,
		rawData.SharesOutstanding, rawData.MarketCap, rawData.EnterpriseValue,
		rawData.WorkingCapital, rawData.CapitalEmployed, rawData.NetDebt,
		rawData.InterestOnLeases, rawData.InterestOnLoans,
		rawData.CompanyType, rawData.NetInterestIncome, rawData.CommissionIncome, rawData.CommissionExpense, rawData.NetCommissionIncome, rawData.CreditLossProvision,
	)

	if err != nil {
		return fmt.Errorf("failed to create metrics: %w", err)
	}

	return nil
}

func (r *RawDataRepository) Update(ctx context.Context, rawData *domain.RawData) error {
	if rawData == nil {
		return fmt.Errorf("rawData is nil: %w", domain.ErrInvalidInput)
	}
	if rawData.Ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if rawData.Year < 1900 || rawData.Year > 2100 {
		return fmt.Errorf("invalid year: %d: %w", rawData.Year, domain.ErrInvalidInput)
	}
	if !rawData.Period.IsValid() {
		return fmt.Errorf("invalid period: %s: %w", rawData.Period, domain.ErrInvalidInput)
	}

	status := rawData.Status
	if !status.IsValid() {
		status = domain.StatusConfirmed
	}

	query := `
		UPDATE metrics SET
			status = $4, report_units = $5,
			revenue = $6, cost_of_revenue = $7, gross_profit = $8, operating_expenses = $9,
			other_income = $10, other_expenses = $11,
			ebit = $12, ebitda = $13, depreciation = $14,
			interest_income = $15, interest_expense = $16,
			profit_before_tax = $17, tax_expense = $18, net_profit = $19, net_profit_parent = $20, basic_eps = $21,
			total_assets = $22, current_assets = $23, cash_and_equivalents = $24, inventories = $25, receivables = $26,
			fixed_assets = $27, right_of_use_assets = $28, intangible_assets = $29, goodwill = $30, total_non_current_assets = $31,
			total_liabilities = $32, current_liabilities = $33, debt = $34, long_term_debt = $35, short_term_debt = $36,
			lt_lease_liabilities = $37, st_lease_liabilities = $38, trade_payables = $39,
			equity = $40, equity_parent = $41, treasury_shares = $42, retained_earnings = $43,
			operating_cash_flow = $44, investing_cash_flow = $45, financing_cash_flow = $46,
			capex = $47, free_cash_flow = $48, dividends_paid = $49, lease_payments = $50,
			acquisitions_net = $51, interest_paid = $52, debt_proceeds = $53, debt_repayments = $54,
			shares_outstanding = $55, market_cap = $56, enterprise_value = $57,
			working_capital = $58, capital_employed = $59, net_debt = $60,
			interest_on_leases = $61, interest_on_loans = $62,
			company_type = $63, net_interest_income = $64, commission_income = $65, commission_expense = $66, net_commission_income = $67, credit_loss_provision = $68,
			updated_at = NOW()
		WHERE ticker = $1 AND year = $2 AND period = $3
	`

	result, err := r.pool.Exec(ctx, query,
		rawData.Ticker, rawData.Year, rawData.Period, status, rawData.ReportUnits,
		rawData.Revenue, rawData.CostOfRevenue, rawData.GrossProfit, rawData.OperatingExpenses,
		rawData.OtherIncome, rawData.OtherExpenses,
		rawData.EBIT, rawData.EBITDA, rawData.Depreciation,
		rawData.InterestIncome, rawData.InterestExpense,
		rawData.ProfitBeforeTax, rawData.TaxExpense, rawData.NetProfit, rawData.NetProfitParent, rawData.BasicEPS,
		rawData.TotalAssets, rawData.CurrentAssets, rawData.CashAndEquivalents, rawData.Inventories, rawData.Receivables,
		rawData.FixedAssets, rawData.RightOfUseAssets, rawData.IntangibleAssets, rawData.Goodwill, rawData.TotalNonCurrentAssets,
		rawData.TotalLiabilities, rawData.CurrentLiabilities, rawData.Debt, rawData.LongTermDebt, rawData.ShortTermDebt,
		rawData.LtLeaseLiabilities, rawData.StLeaseLiabilities, rawData.TradePayables,
		rawData.Equity, rawData.EquityParent, rawData.TreasuryShares, rawData.RetainedEarnings,
		rawData.OperatingCashFlow, rawData.InvestingCashFlow, rawData.FinancingCashFlow,
		rawData.CAPEX, rawData.FreeCashFlow, rawData.DividendsPaid, rawData.LeasePayments,
		rawData.AcquisitionsNet, rawData.InterestPaid, rawData.DebtProceeds, rawData.DebtRepayments,
		rawData.SharesOutstanding, rawData.MarketCap, rawData.EnterpriseValue,
		rawData.WorkingCapital, rawData.CapitalEmployed, rawData.NetDebt,
		rawData.InterestOnLeases, rawData.InterestOnLoans,
		rawData.CompanyType, rawData.NetInterestIncome, rawData.CommissionIncome, rawData.CommissionExpense, rawData.NetCommissionIncome, rawData.CreditLossProvision,
	)

	if err != nil {
		return fmt.Errorf("failed to update metrics: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("metrics not found for ticker %s, year %d, period %s: %w", rawData.Ticker, rawData.Year, rawData.Period, domain.ErrNotFound)
	}

	return nil
}

func (r *RawDataRepository) Delete(ctx context.Context, ticker string, year int, period domain.ReportPeriod) error {
	if ticker == "" {
		return fmt.Errorf("ticker is empty: %w", domain.ErrInvalidInput)
	}
	if year < 1900 || year > 2100 {
		return fmt.Errorf("invalid year: %d: %w", year, domain.ErrInvalidInput)
	}
	if !period.IsValid() {
		return fmt.Errorf("invalid period: %s: %w", period, domain.ErrInvalidInput)
	}

	query := `DELETE FROM metrics WHERE ticker = $1 AND year = $2 AND period = $3`

	result, err := r.pool.Exec(ctx, query, ticker, year, period)
	if err != nil {
		return fmt.Errorf("failed to delete metrics: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("metrics not found for ticker %s, year %d, period %s: %w", ticker, year, period, domain.ErrNotFound)
	}

	return nil
}
