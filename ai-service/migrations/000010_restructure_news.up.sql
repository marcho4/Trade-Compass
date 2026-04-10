ALTER TABLE company_news
    DROP COLUMN IF EXISTS latest_news,
    DROP COLUMN IF EXISTS important_news,
    ADD COLUMN data JSONB NOT NULL DEFAULT '{}';
