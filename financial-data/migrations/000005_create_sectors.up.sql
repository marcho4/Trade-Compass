CREATE TABLE IF NOT EXISTS sectors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для быстрого поиска по имени
CREATE INDEX idx_sectors_name ON sectors(name);

-- Триггер для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_sectors_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_sectors_updated_at
    BEFORE UPDATE ON sectors
    FOR EACH ROW
    EXECUTE FUNCTION update_sectors_updated_at();

-- Вставляем начальные данные из enum Sector
INSERT INTO sectors (id, name) VALUES
    (1, 'Oils'),
    (2, 'Finance'),
    (3, 'Technology'),
    (4, 'Telecom'),
    (5, 'Metals'),
    (6, 'Mining'),
    (7, 'Utilities'),
    (8, 'RealEstate'),
    (9, 'ConsumerStaples'),
    (10, 'ConsumerDiscretionary'),
    (11, 'Healthcare'),
    (12, 'Industrial'),
    (13, 'Energy'),
    (14, 'Materials'),
    (15, 'Transportation'),
    (16, 'Agriculture'),
    (17, 'Chemicals'),
    (18, 'Construction'),
    (19, 'Retail')
ON CONFLICT (id) DO NOTHING;
