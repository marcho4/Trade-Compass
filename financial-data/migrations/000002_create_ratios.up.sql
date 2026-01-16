CREATE TABLE IF NOT EXISTS ratios (
    ticker VARCHAR(10) PRIMARY KEY,
    sector INTEGER NOT NULL,

    price_to_earnings DECIMAL(15, 2),
    price_to_book DECIMAL(15, 2),
    price_to_cash_flow DECIMAL(15, 2),
    ev_to_ebitda DECIMAL(15, 2),
    ev_to_sales DECIMAL(15, 2),
    ev_to_fcf DECIMAL(15, 2),
    peg DECIMAL(15, 2),

    roe DECIMAL(15, 2),
    roa DECIMAL(15, 2),
    roic DECIMAL(15, 2),
    gross_profit_margin DECIMAL(15, 2),
    operating_profit_margin DECIMAL(15, 2),
    net_profit_margin DECIMAL(15, 2),

    current_ratio DECIMAL(15, 2),
    quick_ratio DECIMAL(15, 2),

    net_debt_to_ebitda DECIMAL(15, 2),
    debt_to_equity DECIMAL(15, 2),
    interest_coverage_ratio DECIMAL(15, 2),

    income_quality DECIMAL(15, 2),
    asset_turnover DECIMAL(15, 2),
    inventory_turnover DECIMAL(15, 2),
    receivables_turnover DECIMAL(15, 2),

    eps DECIMAL(15, 2),
    book_value_per_share DECIMAL(15, 2),
    cash_flow_per_share DECIMAL(15, 2),
    dividend_per_share DECIMAL(15, 2),
    dividend_yield DECIMAL(15, 2),
    payout_ratio DECIMAL(15, 2),

    enterprise_value BIGINT,
    market_cap BIGINT,
    free_cash_flow BIGINT,
    capex BIGINT,
    ebitda BIGINT,
    net_debt BIGINT,
    working_capital BIGINT,

    revenue_growth DECIMAL(15, 2),
    earnings_growth DECIMAL(15, 2),
    ebitda_growth DECIMAL(15, 2),
    fcf_growth DECIMAL(15, 2),

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для быстрого поиска по сектору
CREATE INDEX idx_ratios_sector ON ratios(sector);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_ratios_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_ratios_updated_at
    BEFORE UPDATE ON ratios
    FOR EACH ROW
    EXECUTE FUNCTION update_ratios_updated_at();
