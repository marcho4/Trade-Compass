CREATE TABLE IF NOT EXISTS tasks(
    id UUID NOT NULL,
    type VARCHAR(30) NOT NULL,
    pending_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, type)
);