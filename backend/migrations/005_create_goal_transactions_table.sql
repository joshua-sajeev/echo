-- +goose Up

CREATE TABLE goal_transactions (
    id BIGSERIAL PRIMARY KEY,

    goal_id BIGINT NOT NULL
        REFERENCES goals(id)
        ON DELETE CASCADE,

    amount BIGINT NOT NULL
        CHECK (amount <> 0),

    transaction_type VARCHAR(30) NOT NULL
        CHECK (
            transaction_type IN (
                'allocation',
                'manual_contribution',
                'withdrawal'
            )
        ),

    notes TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_goal_transactions_goal_id
ON goal_transactions(goal_id);

CREATE INDEX idx_goal_transactions_created_at
ON goal_transactions(created_at DESC);

-- +goose Down

DROP TABLE IF EXISTS goal_transactions;
