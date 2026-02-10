ALTER TABLE metrics ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'confirmed';
CREATE INDEX idx_metrics_status ON metrics (status);
