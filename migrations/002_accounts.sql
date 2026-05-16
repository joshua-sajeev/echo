-- +goose Up

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,

    name TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    archived_at TIMESTAMP
);

-- +goose Down

DROP TABLE IF EXISTS accounts;
