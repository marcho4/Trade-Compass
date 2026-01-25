DROP TRIGGER IF EXISTS trigger_update_sectors_updated_at ON sectors;
DROP FUNCTION IF EXISTS update_sectors_updated_at();
DROP INDEX IF EXISTS idx_sectors_name;
DROP TABLE IF EXISTS sectors;
