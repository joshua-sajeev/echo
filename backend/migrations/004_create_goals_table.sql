-- +goose Up
CREATE TABLE IF NOT EXISTS goals (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    target_amount BIGINT NOT NULL,
    saved_amount BIGINT NOT NULL DEFAULT 0,
    description VARCHAR(500),
    deadline TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goals_created_at ON goals(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS goals;
