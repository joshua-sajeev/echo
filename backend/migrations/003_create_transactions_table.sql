-- +goose Up
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,

    type TEXT NOT NULL,
    amount NUMERIC NOT NULL,
    name TEXT NOT NULL,

    date DATE NOT NULL DEFAULT CURRENT_DATE,

    from_account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,

    category TEXT,

    jar_id BIGINT REFERENCES jars(id) ON DELETE SET NULL,

    is_master_income BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE transactions;
