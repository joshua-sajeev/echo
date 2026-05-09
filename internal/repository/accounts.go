package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshu-sajeev/echo/internal/models"
)

type AccountRepository interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]models.Account, error)
	ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
}

type pgAccountRepo struct {
	conn *pgx.Conn
}

func NewAccountRepository(conn *pgx.Conn) AccountRepository {
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
		`SELECT id, name, created_at FROM accounts ORDER BY id DESC`,
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

func (r *pgAccountRepo) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
			a.id,
			a.name,
			a.created_at,
			COALESCE(SUM(CASE WHEN t.to_account_id   = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON t.from_account_id = a.id OR t.to_account_id = a.id
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
