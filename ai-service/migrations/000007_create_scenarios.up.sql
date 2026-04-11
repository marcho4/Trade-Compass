CREATE TABLE IF NOT EXISTS scenarios (
    id VARCHAR(50) NOT NULL,
    ticker VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    probability DOUBLE PRECISION NOT NULL,
    terminal_growth_rate DOUBLE PRECISION NOT NULL,
    growth_factors_applied JSONB,
    risks_applied JSONB,
    assumptions JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, ticker)
);
