CREATE TABLE IF NOT EXISTS company_profiles (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(20) NOT NULL UNIQUE,
    company_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    products_and_services JSONB NOT NULL DEFAULT '[]',
    markets JSONB NOT NULL DEFAULT '[]',
    key_clients TEXT NOT NULL DEFAULT '',
    business_model TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS revenue_sources (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(20) NOT NULL,
    segment VARCHAR(255) NOT NULL,
    share_pct REAL NOT NULL DEFAULT 0,
    approximate BOOLEAN NOT NULL DEFAULT false,
    description TEXT NOT NULL DEFAULT '',
    trend VARCHAR(20) NOT NULL DEFAULT 'stable',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(ticker, segment)
);

CREATE TABLE IF NOT EXISTS company_dependencies (
    id BIGSERIAL PRIMARY KEY,
    ticker VARCHAR(20) NOT NULL,
    factor VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(ticker, factor)
);
