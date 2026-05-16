-- +goose Up

CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,

    type TEXT NOT NULL,
    -- income | expense | transfer

    amount NUMERIC NOT NULL,

    name TEXT NOT NULL,

    date DATE NOT NULL DEFAULT CURRENT_DATE,

    from_account_id BIGINT
        REFERENCES accounts(id)
        ON DELETE SET NULL,

    to_account_id BIGINT
        REFERENCES accounts(id)
        ON DELETE SET NULL,

    category TEXT,

    target_jar_id BIGINT
        REFERENCES jars(id)
        ON DELETE SET NULL,

    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

-- +goose Down

DROP TABLE IF EXISTS transactions;
