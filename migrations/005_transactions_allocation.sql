-- +goose Up

CREATE TABLE transaction_allocations (
    id BIGSERIAL PRIMARY KEY,

    transaction_id BIGINT NOT NULL
        REFERENCES transactions(id)
        ON DELETE CASCADE,

    jar_id BIGINT NOT NULL
        REFERENCES jars(id)
        ON DELETE CASCADE,

    amount NUMERIC NOT NULL,

    created_at TIMESTAMP DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS transaction_allocations;
