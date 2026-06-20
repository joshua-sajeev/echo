-- +goose Up

CREATE TABLE IF NOT EXISTS goals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,

    target_amount BIGINT NOT NULL CHECK (target_amount > 0),
    saved_amount BIGINT NOT NULL DEFAULT 0 CHECK (saved_amount >= 0),

    deadline TIMESTAMP,

    allocation_percentage INT NOT NULL DEFAULT 0
        CHECK (
            allocation_percentage >= 0
            AND allocation_percentage <= 100
        ),

    is_archived BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_goals_created_at
ON goals(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_goals_archived
ON goals(is_archived);

-- +goose Down

DROP TABLE IF EXISTS goals;
