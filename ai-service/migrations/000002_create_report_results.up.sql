CREATE TABLE IF NOT EXISTS report_results (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    year INTEGER NOT NULL,
    period INTEGER NOT NULL,
    health INTEGER NOT NULL,
    growth INTEGER NOT NULL,
    moat INTEGER NOT NULL,
    dividends INTEGER NOT NULL,
    value INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_report_results_key ON report_results(ticker, year, period);

CREATE OR REPLACE FUNCTION update_report_results_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_report_results_updated_at
    BEFORE UPDATE ON report_results
    FOR EACH ROW
    EXECUTE FUNCTION update_report_results_updated_at();
