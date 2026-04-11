ALTER TABLE scenarios DROP CONSTRAINT IF EXISTS scenarios_pkey;
ALTER TABLE scenarios DROP COLUMN IF EXISTS task_id;
ALTER TABLE scenarios ADD PRIMARY KEY (id, ticker);
