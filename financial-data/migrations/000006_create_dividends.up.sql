CREATE TABLE IF NOT EXISTS dividends (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    ex_dividend_date DATE NOT NULL,
    payment_date DATE NOT NULL,
    amount_per_share DECIMAL(10, 2) NOT NULL,
    dividend_yield DECIMAL(5, 2),
    payout_ratio DECIMAL(5, 2),
    currency VARCHAR(3) NOT NULL DEFAULT 'RUB',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_dividends_ticker ON dividends(ticker);
CREATE INDEX idx_dividends_ticker_date ON dividends(ticker, ex_dividend_date DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_dividends_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_dividends_updated_at
    BEFORE UPDATE ON dividends
    FOR EACH ROW
    EXECUTE FUNCTION update_dividends_updated_at();
