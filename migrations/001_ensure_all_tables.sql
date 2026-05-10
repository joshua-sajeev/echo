-- migrations/001_ensure_all_tables.sql
-- Idempotent: safe to run multiple times, on a fresh DB or an existing one.
-- Run with: psql $DATABASE_URL -f migrations/001_ensure_all_tables.sql

-- ── accounts ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS accounts (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT      NOT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    archived_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE accounts
    ADD COLUMN IF NOT EXISTS archived_at TIMESTAMP DEFAULT NULL;

-- ── jars ─────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS jars (
    id               BIGSERIAL PRIMARY KEY,
    name             TEXT      NOT NULL UNIQUE,
    created_at       TIMESTAMP DEFAULT now(),
    allocation_type  TEXT      NOT NULL DEFAULT 'fixed',
    allocation_value NUMERIC   NOT NULL DEFAULT 0,
    round_to         INTEGER   NOT NULL DEFAULT 0,
    sort_order       INTEGER   NOT NULL DEFAULT 0,
    is_system        BOOLEAN   NOT NULL DEFAULT false
);

ALTER TABLE jars
    ADD COLUMN IF NOT EXISTS allocation_type  TEXT    NOT NULL DEFAULT 'fixed',
    ADD COLUMN IF NOT EXISTS allocation_value NUMERIC NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS round_to         INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS sort_order       INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_system        BOOLEAN NOT NULL DEFAULT false;

-- ── transactions ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS transactions (
    id              BIGSERIAL PRIMARY KEY,
    type            TEXT      NOT NULL,
    amount          NUMERIC   NOT NULL,
    name            TEXT      NOT NULL,
    date            DATE      NOT NULL DEFAULT now(),
    from_account_id BIGINT    REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id   BIGINT    REFERENCES accounts(id) ON DELETE SET NULL,
    category        TEXT,
    sub_category    TEXT,
    jar_id          BIGINT    REFERENCES jars(id) ON DELETE SET NULL,
    is_master_income BOOLEAN  DEFAULT false,
    created_at      TIMESTAMP DEFAULT now()
);

ALTER TABLE transactions
    ADD COLUMN IF NOT EXISTS from_account_id  BIGINT  REFERENCES accounts(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS to_account_id    BIGINT  REFERENCES accounts(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS category         TEXT,
    ADD COLUMN IF NOT EXISTS sub_category     TEXT,
    ADD COLUMN IF NOT EXISTS jar_id           BIGINT  REFERENCES jars(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS is_master_income BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS created_at       TIMESTAMP DEFAULT now();

-- ── tx_templates ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS tx_templates (
    id               BIGSERIAL PRIMARY KEY,
    name             TEXT      NOT NULL,
    type             TEXT      NOT NULL DEFAULT 'expense',
    amount           NUMERIC   NOT NULL DEFAULT 0,
    jar_id           BIGINT    REFERENCES jars(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    from_account_id  BIGINT    REFERENCES accounts(id) ON DELETE SET NULL,
    to_account_id    BIGINT    REFERENCES accounts(id) ON DELETE SET NULL,
    is_master_income BOOLEAN   NOT NULL DEFAULT false
);

ALTER TABLE tx_templates
    ADD COLUMN IF NOT EXISTS from_account_id  BIGINT  REFERENCES accounts(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS to_account_id    BIGINT  REFERENCES accounts(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS is_master_income BOOLEAN NOT NULL DEFAULT false;

-- ── default system jars (idempotent via ON CONFLICT DO NOTHING) ───────────────
INSERT INTO jars (name, allocation_type, allocation_value, round_to, sort_order, is_system) VALUES
    ('Charity',      'percent_total',     10,   100, 1, true),
    ('SIP',          'fixed',           1000,     0, 2, true),
    ('Chitty',       'cap',             5000,     0, 3, true),
    ('Necessities',  'percent_remainder',  0,   500, 4, true),
    ('Leisure',      'remainder',          0,   100, 5, true)
ON CONFLICT (name) DO NOTHING;
