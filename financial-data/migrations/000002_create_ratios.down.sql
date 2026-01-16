DROP TRIGGER IF EXISTS trigger_update_ratios_updated_at ON ratios;
DROP FUNCTION IF EXISTS update_ratios_updated_at();
DROP INDEX IF EXISTS idx_ratios_sector;
DROP TABLE IF EXISTS ratios;
