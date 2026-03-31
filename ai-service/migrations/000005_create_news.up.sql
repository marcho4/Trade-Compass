CREATE TABLE IF NOT EXISTS company_news(
    ticker VARCHAR(20) PRIMARY KEY,
    latest_news JSONB NOT NULL,
    important_news JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);