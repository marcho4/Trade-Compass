DROP TRIGGER IF EXISTS trigger_update_dividends_updated_at ON dividends;
DROP FUNCTION IF EXISTS update_dividends_updated_at();
DROP INDEX IF EXISTS idx_dividends_ticker_date;
DROP INDEX IF EXISTS idx_dividends_ticker;
DROP TABLE IF EXISTS dividends;
