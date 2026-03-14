CREATE TABLE IF NOT EXISTS subscription_types (
    id INT PRIMARY KEY,
    name VARCHAR(20) UNIQUE NOT NULL
);

INSERT INTO subscription_types (id, name) VALUES 
(1, 'free'), 
(2, 'pro'), 
(3, 'premium');

CREATE TABLE IF NOT EXISTS subscriptions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_type_id INT REFERENCES subscription_types(id) DEFAULT 1,
    start_date TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_subscription_type ON subscriptions(subscription_type_id);
