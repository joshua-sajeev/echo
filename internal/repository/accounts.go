// Package repository
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type AccountRepository interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]models.Account, error)
	ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	Rename(ctx context.Context, id int64, name string) error
	Archive(ctx context.Context, id int64) error
	Unarchive(ctx context.Context, id int64) error
}

type pgAccountRepo struct {
	conn *pgxpool.Pool
}

func NewAccountRepository(conn *pgxpool.Pool) AccountRepository {
	return &pgAccountRepo{conn: conn}
}

func (r *pgAccountRepo) Create(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.conn.QueryRow(ctx,
		`INSERT INTO accounts (name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&id)
	return id, err
}

func (r *pgAccountRepo) List(ctx context.Context) ([]models.Account, error) {
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

func (r *pgAccountRepo) listWithBalancesWhere(ctx context.Context, where string) ([]models.AccountWithBalance, error) {
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

func (r *pgAccountRepo) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, "a.archived_at IS NULL")
}

func (r *pgAccountRepo) ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, "a.archived_at IS NOT NULL")
}

func (r *pgAccountRepo) Rename(ctx context.Context, id int64, name string) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET name = $1 WHERE id = $2 AND archived_at IS NULL`,
		name, id,
	)
	return err
}

func (r *pgAccountRepo) Archive(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET archived_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}

func (r *pgAccountRepo) Unarchive(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx,
		`UPDATE accounts SET archived_at = NULL WHERE id = $1`,
		id,
	)
	return err
}
