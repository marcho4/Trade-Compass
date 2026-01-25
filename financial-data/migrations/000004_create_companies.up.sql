CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    inn VARCHAR(12) UNIQUE NOT NULL,
    ticker VARCHAR(10) UNIQUE NOT NULL,
    owner TEXT NOT NULL,
    sector_id INTEGER NOT NULL,
    lot_size INTEGER NOT NULL,
    ceo VARCHAR(100),
    employees INTEGER,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_companies_ticker ON companies(ticker);
CREATE INDEX idx_companies_sector ON companies(sector_id);
CREATE INDEX idx_companies_inn ON companies(inn);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_companies_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_companies_updated_at
    BEFORE UPDATE ON companies
    FOR EACH ROW
    EXECUTE FUNCTION update_companies_updated_at();
