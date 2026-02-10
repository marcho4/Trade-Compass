DROP INDEX IF EXISTS idx_metrics_status;
ALTER TABLE metrics DROP COLUMN IF EXISTS status;
