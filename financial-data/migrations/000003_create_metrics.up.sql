CREATE TABLE IF NOT EXISTS metrics (
    ticker VARCHAR(10) NOT NULL,
    year INTEGER NOT NULL,
    period VARCHAR(10) NOT NULL,

    -- P&L (Отчёт о прибылях и убытках)
    revenue BIGINT,
    cost_of_revenue BIGINT,
    gross_profit BIGINT,
    operating_expenses BIGINT,
    ebit BIGINT,
    ebitda BIGINT,
    interest_expense BIGINT,
    tax_expense BIGINT,
    net_profit BIGINT,

    -- Balance Sheet (Баланс)
    total_assets BIGINT,
    current_assets BIGINT,
    cash_and_equivalents BIGINT,
    inventories BIGINT,
    receivables BIGINT,

    total_liabilities BIGINT,
    current_liabilities BIGINT,
    debt BIGINT,
    long_term_debt BIGINT,
    short_term_debt BIGINT,
    equity BIGINT,
    retained_earnings BIGINT,

    -- Cash Flow Statement (Отчёт о движении денежных средств)
    operating_cash_flow BIGINT,
    investing_cash_flow BIGINT,
    financing_cash_flow BIGINT,
    capex BIGINT,
    free_cash_flow BIGINT,

    -- Market Data (для мультипликаторов)
    shares_outstanding BIGINT,
    market_cap BIGINT,

    -- Calculated fields (Дополнительные расчётные поля)
    working_capital BIGINT,
    capital_employed BIGINT,
    enterprise_value BIGINT,
    net_debt BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (ticker, year, period)
);

-- Индекс для быстрого поиска по тикеру
CREATE INDEX idx_metrics_ticker ON metrics(ticker);

-- Индекс для поиска последних метрик
CREATE INDEX idx_metrics_ticker_year_period ON metrics(ticker, year DESC, period);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_metrics_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_metrics_updated_at
    BEFORE UPDATE ON metrics
    FOR EACH ROW
    EXECUTE FUNCTION update_metrics_updated_at();
