-- +goose Up

CREATE TABLE jars (
    id BIGSERIAL PRIMARY KEY,

    name TEXT NOT NULL UNIQUE,

    allocation_type TEXT NOT NULL DEFAULT 'fixed',
    -- fixed | percent | remainder | capped

    allocation_value NUMERIC NOT NULL DEFAULT 0,
    -- meaning depends on allocation_type

    round_to INTEGER NOT NULL DEFAULT 0,
    -- rounding rule (100, 500 etc.)

    sort_order INTEGER NOT NULL DEFAULT 0,

    is_active BOOLEAN DEFAULT true,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS jars;
