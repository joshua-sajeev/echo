-- +goose Up

-- Transactions

CREATE INDEX idx_transactions_date
ON transactions(date DESC);

CREATE INDEX idx_transactions_type
ON transactions(type);

CREATE INDEX idx_transactions_from_account_id
ON transactions(from_account_id);

CREATE INDEX idx_transactions_to_account_id
ON transactions(to_account_id);

CREATE INDEX idx_transactions_target_jar_id
ON transactions(target_jar_id);

CREATE INDEX idx_transactions_category
ON transactions(category);


-- Transaction Allocations

CREATE INDEX idx_transaction_allocations_transaction_id
ON transaction_allocations(transaction_id);

CREATE INDEX idx_transaction_allocations_jar_id
ON transaction_allocations(jar_id);


-- Jars

CREATE INDEX idx_jars_is_active
ON jars(is_active);

CREATE INDEX idx_jars_sort_order
ON jars(sort_order);


-- Accounts

CREATE INDEX idx_accounts_archived_at
ON accounts(archived_at);


-- +goose Down

DROP INDEX IF EXISTS idx_transactions_date;

DROP INDEX IF EXISTS idx_transactions_type;

DROP INDEX IF EXISTS idx_transactions_from_account_id;

DROP INDEX IF EXISTS idx_transactions_to_account_id;

DROP INDEX IF EXISTS idx_transactions_target_jar_id;

DROP INDEX IF EXISTS idx_transactions_category;

DROP INDEX IF EXISTS idx_transaction_allocations_transaction_id;

DROP INDEX IF EXISTS idx_transaction_allocations_jar_id;

DROP INDEX IF EXISTS idx_jars_is_active;

DROP INDEX IF EXISTS idx_jars_sort_order;

DROP INDEX IF EXISTS idx_accounts_archived_at;
