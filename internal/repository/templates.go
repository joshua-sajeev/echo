package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type TxTemplateRepository interface {
	List(ctx context.Context) ([]models.TxTemplate, error)
	Create(ctx context.Context, t models.TxTemplate) (int64, error)
	Delete(ctx context.Context, id int64) error
	EnsureTable(ctx context.Context) error
}

type pgTxTemplateRepo struct {
	conn *pgxpool.Pool
}

func NewTxTemplateRepository(conn *pgxpool.Pool) TxTemplateRepository {
	return &pgTxTemplateRepo{conn: conn}
}

// EnsureTable creates the table if it doesn't exist and adds new columns
// (from_account_id, to_account_id, is_master_income) if they're missing.
// Safe to call every startup.
func (r *pgTxTemplateRepo) EnsureTable(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS tx_templates (
			id               BIGSERIAL PRIMARY KEY,
			name             TEXT NOT NULL,
			type             TEXT NOT NULL DEFAULT 'expense',
			amount           NUMERIC(12,2) NOT NULL DEFAULT 0,
			jar_id           BIGINT REFERENCES jars(id) ON DELETE SET NULL,
			from_account_id  BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
			to_account_id    BIGINT REFERENCES accounts(id) ON DELETE SET NULL,
			is_master_income BOOLEAN NOT NULL DEFAULT false,
			created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Idempotently add columns for existing tables that predate this version
	for _, col := range []struct{ name, def string }{
		{"from_account_id", "BIGINT REFERENCES accounts(id) ON DELETE SET NULL"},
		{"to_account_id", "BIGINT REFERENCES accounts(id) ON DELETE SET NULL"},
		{"is_master_income", "BOOLEAN NOT NULL DEFAULT false"},
	} {
		_, _ = r.conn.Exec(ctx,
			`ALTER TABLE tx_templates ADD COLUMN IF NOT EXISTS `+col.name+` `+col.def,
		)
	}

	return nil
}

func (r *pgTxTemplateRepo) List(ctx context.Context) ([]models.TxTemplate, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
			t.id, t.name, t.type, t.amount,
			t.jar_id,           COALESCE(j.name,  ''),
			t.from_account_id,  COALESCE(fa.name, ''),
			t.to_account_id,    COALESCE(ta.name, ''),
			t.is_master_income,
			t.created_at
		FROM tx_templates t
		LEFT JOIN jars     j  ON j.id  = t.jar_id
		LEFT JOIN accounts fa ON fa.id = t.from_account_id
		LEFT JOIN accounts ta ON ta.id = t.to_account_id
		ORDER BY t.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.TxTemplate
	for rows.Next() {
		var t models.TxTemplate
		if err := rows.Scan(
			&t.ID, &t.Name, &t.Type, &t.Amount,
			&t.JarID, &t.JarName,
			&t.FromAccountID, &t.FromAccountName,
			&t.ToAccountID, &t.ToAccountName,
			&t.IsMasterIncome,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *pgTxTemplateRepo) Create(ctx context.Context, t models.TxTemplate) (int64, error) {
	var id int64
	err := r.conn.QueryRow(ctx, `
		INSERT INTO tx_templates
			(name, type, amount, jar_id, from_account_id, to_account_id, is_master_income)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`, t.Name, t.Type, t.Amount, t.JarID,
		t.FromAccountID, t.ToAccountID, t.IsMasterIncome,
	).Scan(&id)
	return id, err
}

func (r *pgTxTemplateRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx, `DELETE FROM tx_templates WHERE id = $1`, id)
	return err
}
