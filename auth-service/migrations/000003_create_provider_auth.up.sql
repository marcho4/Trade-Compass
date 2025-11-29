CREATE TABLE provider_auth (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_type VARCHAR(20) NOT NULL CHECK (provider_type IN ('google', 'yandex')),
    provider_user_id VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    
    UNIQUE(provider_type, provider_user_id)
);

CREATE INDEX idx_provider_auth_user_id ON provider_auth(user_id);
