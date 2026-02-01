DROP TRIGGER IF EXISTS trigger_update_companies_updated_at ON companies;
DROP FUNCTION IF EXISTS update_companies_updated_at();
DROP INDEX IF EXISTS idx_companies_inn;
DROP INDEX IF EXISTS idx_companies_sector;
DROP INDEX IF EXISTS idx_companies_ticker;
DROP TABLE IF EXISTS companies;
