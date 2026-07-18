-- +goose Up
CREATE TABLE jars (
    id BIGSERIAL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,

    allocation_type TEXT NOT NULL CHECK (
        allocation_type IN (
            'percentage',
            'remainder'
        )
    ),

    value BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE jars;
