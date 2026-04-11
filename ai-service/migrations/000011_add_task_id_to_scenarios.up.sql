ALTER TABLE scenarios ADD COLUMN IF NOT EXISTS task_id VARCHAR(255);

UPDATE scenarios SET task_id = id WHERE task_id IS NULL;

ALTER TABLE scenarios ALTER COLUMN task_id SET NOT NULL;

ALTER TABLE scenarios DROP CONSTRAINT IF EXISTS scenarios_pkey;

ALTER TABLE scenarios ADD PRIMARY KEY (task_id, ticker, id);

CREATE INDEX IF NOT EXISTS idx_scenarios_task_id_ticker ON scenarios(task_id, ticker);
