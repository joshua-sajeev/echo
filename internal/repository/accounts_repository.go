// Package repository provides database access implementations
// for account-related operations.
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

// AccountRepository handles account database operations.
type AccountRepository struct {
	conn *pgxpool.Pool
}

// NewAccountRepository creates a new AccountRepository instance.
func NewAccountRepository(conn *pgxpool.Pool) AccountRepositoryInterface {
	return &AccountRepository{conn: conn}
}

// Create inserts a new account into the database
// and returns the generated account ID.
func (r *AccountRepository) Create(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.conn.QueryRow(ctx,
		`INSERT INTO accounts (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)
	return id, err
}

// List returns all non-archived accounts ordered by newest first.
func (r *AccountRepository) List(ctx context.Context) ([]models.Account, error) {
	rows, err := r.conn.Query(ctx,
		`SELECT id, name, created_at FROM accounts WHERE archived_at IS NULL ORDER BY id DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

// listWithBalancesWhere returns accounts with computed balances
// using the provided SQL WHERE condition.
func (r *AccountRepository) listWithBalancesWhere(ctx context.Context, where string) ([]models.AccountWithBalance, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
			a.id, a.name, a.created_at,
			COALESCE(SUM(CASE WHEN t.to_account_id   = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON t.from_account_id = a.id OR t.to_account_id = a.id
		WHERE `+where+`
		GROUP BY a.id, a.name, a.created_at
		ORDER BY a.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.AccountWithBalance

	for rows.Next() {
		var a models.AccountWithBalance

		if err := rows.Scan(&a.ID, &a.Name, &a.CreatedAt, &a.Balance); err != nil {
			return nil, err
		}

		accounts = append(accounts, a)
	}

	return accounts, rows.Err()
}

// ListWithBalances returns all active accounts
// along with their calculated balances.
func (r *AccountRepository) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, "a.archived_at IS NULL")
}

// ListArchivedWithBalances returns all archived accounts
// along with their calculated balances.
func (r *AccountRepository) ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, "a.archived_at IS NOT NULL")
}

// Rename updates the name of an active account.
func (r *AccountRepository) Rename(ctx context.Context, id int64, name string) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET name = $1 WHERE id = $2 AND archived_at IS NULL`,
		name, id,
	)

	return err
}

// Archive marks an account as archived.
func (r *AccountRepository) Archive(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET archived_at = NOW() WHERE id = $1`,
		id,
	)

	return err
}

// Unarchive restores an archived account.
func (r *AccountRepository) Unarchive(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET archived_at = NULL WHERE id = $1`,
		id,
	)

	return err
}
