// Package repository handles repository functions for the models
package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type AccountRepo struct {
	conn *pgxpool.Pool
}

type AccountRepositoryInterface interface {
	Create(ctx context.Context, name string) (int64, error)
	List(ctx context.Context) ([]models.Account, error)
	ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error)
	Rename(ctx context.Context, id int64, name string) error
	Archive(ctx context.Context, id int64) error
	Unarchive(ctx context.Context, id int64) error
}

var _ AccountRepositoryInterface = (*AccountRepo)(nil)

// NewAccountRepository creates a new account repository
func NewAccountRepository(conn *pgxpool.Pool) *AccountRepo {
	return &AccountRepo{conn: conn}
}

// Create inserts a new account and returns its ID
func (r *AccountRepo) Create(ctx context.Context, name string) (int64, error) {
	var id int64

	err := r.conn.QueryRow(
		ctx,
		`INSERT INTO accounts (name, is_archived)
		 VALUES ($1, false)
		 RETURNING id`,
		name,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create account: %w", err)
	}

	return id, nil
}

// List returns all non-archived accounts
func (r *AccountRepo) List(ctx context.Context) ([]models.Account, error) {
	rows, err := r.conn.Query(ctx,
		`SELECT id, name, is_archived, created_at FROM accounts WHERE is_archived = false ORDER BY id DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var a models.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.IsArchived, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accounts: %w", err)
	}

	return accounts, nil
}

// listWithBalancesWhere is a helper function to query accounts with balances
// archived parameter: true for archived accounts, false for active accounts
func (r *AccountRepo) listWithBalancesWhere(ctx context.Context, archived bool) ([]models.AccountWithBalance, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT
			a.id, 
			a.name, 
			a.is_archived,
			a.created_at,
			COALESCE(SUM(CASE WHEN t.to_account_id = a.id THEN t.amount ELSE 0 END), 0) -
			COALESCE(SUM(CASE WHEN t.from_account_id = a.id THEN t.amount ELSE 0 END), 0) AS balance
		FROM accounts a
		LEFT JOIN transactions t ON t.from_account_id = a.id OR t.to_account_id = a.id
		WHERE a.is_archived = $1
		GROUP BY a.id, a.name, a.is_archived, a.created_at
		ORDER BY a.id DESC
	`, archived)
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts with balances: %w", err)
	}
	defer rows.Close()

	var accounts []models.AccountWithBalance
	for rows.Next() {
		var a models.AccountWithBalance
		if err := rows.Scan(&a.ID, &a.Name, &a.IsArchived, &a.CreatedAt, &a.Balance); err != nil {
			return nil, fmt.Errorf("failed to scan account with balance: %w", err)
		}
		accounts = append(accounts, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accounts with balances: %w", err)
	}

	return accounts, nil
}

// ListWithBalances returns all non-archived accounts with their balances
func (r *AccountRepo) ListWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, false)
}

// ListArchivedWithBalances returns all archived accounts with their balances
func (r *AccountRepo) ListArchivedWithBalances(ctx context.Context) ([]models.AccountWithBalance, error) {
	return r.listWithBalancesWhere(ctx, true)
}

// Rename updates an account name (only if not archived)
func (r *AccountRepo) Rename(ctx context.Context, id int64, name string) error {
	if id <= 0 {
		return fmt.Errorf("invalid account id: %d", id)
	}
	if name == "" {
		return fmt.Errorf("account name cannot be empty")
	}

	tag, err := r.conn.Exec(ctx,
		`UPDATE accounts SET name = $1 WHERE id = $2 AND is_archived = false`,
		name, id,
	)
	if err != nil {
		return fmt.Errorf("failed to rename account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account not found or is archived")
	}

	return nil
}

// Archive marks an account as archived
func (r *AccountRepo) Archive(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid account id: %d", id)
	}

	tag, err := r.conn.Exec(ctx,
		`UPDATE accounts SET is_archived = true WHERE id = $1 AND is_archived = false`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to archive account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account not found or already archived")
	}

	return nil
}

// Unarchive marks an account as active
func (r *AccountRepo) Unarchive(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid account id: %d", id)
	}

	tag, err := r.conn.Exec(ctx,
		`UPDATE accounts SET is_archived = false WHERE id = $1 AND is_archived = true`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to unarchive account: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("account not found or not archived")
	}

	return nil
}
