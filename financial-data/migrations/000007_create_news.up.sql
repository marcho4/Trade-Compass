CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(10),
    sector_id INTEGER,
    date TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    source VARCHAR(100) NOT NULL,
    url TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индексы для быстрого поиска
CREATE INDEX idx_news_ticker ON news(ticker);
CREATE INDEX idx_news_sector ON news(sector_id);
CREATE INDEX idx_news_date ON news(date DESC);
CREATE INDEX idx_news_ticker_date ON news(ticker, date DESC);
CREATE INDEX idx_news_sector_date ON news(sector_id, date DESC);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_news_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_news_updated_at
    BEFORE UPDATE ON news
    FOR EACH ROW
    EXECUTE FUNCTION update_news_updated_at();
