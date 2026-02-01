DROP TRIGGER IF EXISTS trigger_update_news_updated_at ON news;
DROP FUNCTION IF EXISTS update_news_updated_at();
DROP INDEX IF EXISTS idx_news_sector_date;
DROP INDEX IF EXISTS idx_news_ticker_date;
DROP INDEX IF EXISTS idx_news_date;
DROP INDEX IF EXISTS idx_news_sector;
DROP INDEX IF EXISTS idx_news_ticker;
DROP TABLE IF EXISTS news;
