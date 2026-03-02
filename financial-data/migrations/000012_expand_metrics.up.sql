ALTER TABLE metrics
    ADD COLUMN report_units VARCHAR(20),

    -- Income Statement
    ADD COLUMN other_income BIGINT,
    ADD COLUMN other_expenses BIGINT,
    ADD COLUMN depreciation BIGINT,
    ADD COLUMN interest_income BIGINT,
    ADD COLUMN profit_before_tax BIGINT,
    ADD COLUMN net_profit_parent BIGINT,
    ADD COLUMN basic_eps DOUBLE PRECISION,

    -- Balance Sheet
    ADD COLUMN fixed_assets BIGINT,
    ADD COLUMN right_of_use_assets BIGINT,
    ADD COLUMN intangible_assets BIGINT,
    ADD COLUMN goodwill BIGINT,
    ADD COLUMN total_non_current_assets BIGINT,
    ADD COLUMN lt_lease_liabilities BIGINT,
    ADD COLUMN st_lease_liabilities BIGINT,
    ADD COLUMN trade_payables BIGINT,
    ADD COLUMN equity_parent BIGINT,
    ADD COLUMN treasury_shares BIGINT,

    -- Cash Flow
    ADD COLUMN dividends_paid BIGINT,
    ADD COLUMN lease_payments BIGINT,
    ADD COLUMN acquisitions_net BIGINT,
    ADD COLUMN interest_paid BIGINT,
    ADD COLUMN debt_proceeds BIGINT,
    ADD COLUMN debt_repayments BIGINT,

    -- Notes Breakdown
    ADD COLUMN interest_on_leases BIGINT,
    ADD COLUMN interest_on_loans BIGINT;
