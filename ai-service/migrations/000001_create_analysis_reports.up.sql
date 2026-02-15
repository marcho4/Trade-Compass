CREATE TABLE IF NOT EXISTS analysis_reports (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    year INTEGER NOT NULL,
    period INTEGER NOT NULL,
    analysis TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_analysis_reports_key ON analysis_reports(ticker, year, period);

CREATE OR REPLACE FUNCTION update_analysis_reports_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_analysis_reports_updated_at
    BEFORE UPDATE ON analysis_reports
    FOR EACH ROW
    EXECUTE FUNCTION update_analysis_reports_updated_at();
