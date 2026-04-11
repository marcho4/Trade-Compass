CREATE TABLE IF NOT EXISTS risk_and_growth (
    ticker VARCHAR(20) PRIMARY KEY,
    factors JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
