DROP TRIGGER IF EXISTS trigger_update_metrics_updated_at ON metrics;
DROP FUNCTION IF EXISTS update_metrics_updated_at();
DROP INDEX IF EXISTS idx_metrics_ticker_year_period;
DROP INDEX IF EXISTS idx_metrics_ticker;
DROP TABLE IF EXISTS metrics;
