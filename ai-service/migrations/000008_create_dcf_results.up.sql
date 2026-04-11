CREATE TABLE IF NOT EXISTS dcf_results (
    id VARCHAR(50) NOT NULL,
    ticker VARCHAR(20) NOT NULL,
    scenario_type VARCHAR(50) NOT NULL,
    probability DOUBLE PRECISION NOT NULL,
    enterprise_value DOUBLE PRECISION NOT NULL,
    equity_value DOUBLE PRECISION NOT NULL,
    price_per_share DOUBLE PRECISION NOT NULL,
    terminal_value DOUBLE PRECISION NOT NULL,
    yearly_fcfs JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, ticker, scenario_type)
);
