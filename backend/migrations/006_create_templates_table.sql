-- +goose Up
CREATE TABLE IF NOT EXISTS transaction_templates (
    id BIGSERIAL PRIMARY KEY,
    template_name TEXT NOT NULL,
    type TEXT NOT NULL,
    amount BIGINT,
    name TEXT NOT NULL,
    category TEXT,
    account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
    from_account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
    jar_id BIGINT REFERENCES jars(id) ON DELETE SET NULL,
    is_master_income BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS transaction_templates;
