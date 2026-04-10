ALTER TABLE company_news
    DROP COLUMN IF EXISTS data,
    ADD COLUMN latest_news JSONB NOT NULL DEFAULT '[]',
    ADD COLUMN important_news JSONB NOT NULL DEFAULT '[]';
