package transactions

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	db *pgxpool.Pool
}

type TransactionRepositoryInterface interface {
	Create(ctx context.Context, tx Transaction) (int64, error)
	List(ctx context.Context) ([]Transaction, error)
	GetByID(ctx context.Context, id int64) (*Transaction, error)
	Update(ctx context.Context, tx Transaction) error
	Delete(ctx context.Context, id int64) error
	GetCurrentMonthIncome(ctx context.Context) (int64, error)
}
type TransactionIncomeRepository interface{}

var _ TransactionRepositoryInterface = (*TransactionRepository)(nil)

func NewTransactionRepository(db *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx Transaction) (int64, error) {
	var id int64

	err := r.db.QueryRow(ctx, `
		INSERT INTO transactions (
			type,
			amount,
			name,
			date,
			from_account_id,
			to_account_id,
			category,
			jar_id,
			is_master_income
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id
	`,
		tx.Type,
		tx.Amount,
		tx.Name,
		tx.Date,
		tx.FromAccountID,
		tx.ToAccountID,
		tx.Category,
		tx.JarID,
		tx.IsMasterIncome,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create transaction: %w", err)
	}

	return id, nil
}

func (r *TransactionRepository) List(ctx context.Context) ([]Transaction, error) {
	result := make([]Transaction, 0)

	rows, err := r.db.Query(ctx, `
    SELECT
        id,
        type,
        amount,
        name,
        date,
        from_account_id,
        to_account_id,
        category,
        jar_id,
        is_master_income,
        created_at
    FROM transactions
    ORDER BY date DESC
`)
	if err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction

		err := rows.Scan(
			&tx.ID,
			&tx.Type,
			&tx.Amount,
			&tx.Name,
			&tx.Date,
			&tx.FromAccountID,
			&tx.ToAccountID,
			&tx.Category,
			&tx.JarID,
			&tx.IsMasterIncome,
			&tx.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan transaction: %w", err)
		}

		result = append(result, tx)
	}

	return result, nil
}

func (r *TransactionRepository) GetByID(ctx context.Context, id int64) (*Transaction, error) {
	var tx Transaction
	err := r.db.QueryRow(ctx, `
		SELECT
			id, type, amount, name, date,
			from_account_id, to_account_id,
			category, jar_id, is_master_income, created_at
		FROM transactions
		WHERE id = $1
	`, id).Scan(
		&tx.ID,
		&tx.Type,
		&tx.Amount,
		&tx.Name,
		&tx.Date,
		&tx.FromAccountID,
		&tx.ToAccountID,
		&tx.Category,
		&tx.JarID,
		&tx.IsMasterIncome,
		&tx.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("get transaction by id: %w", err)
	}
	return &tx, nil
}

func (r *TransactionRepository) Update(ctx context.Context, tx Transaction) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE transactions
		SET
			type = $1,
			amount = $2,
			name = $3,
			date = $4,
			from_account_id = $5,
			to_account_id = $6,
			category = $7,
			jar_id = $8,
			is_master_income = $9
		WHERE id = $10
	`,
		tx.Type,
		tx.Amount,
		tx.Name,
		tx.Date,
		tx.FromAccountID,
		tx.ToAccountID,
		tx.Category,
		tx.JarID,
		tx.IsMasterIncome,
		tx.ID,
	)
	if err != nil {
		return fmt.Errorf("update transaction: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

func (r *TransactionRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM transactions WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("delete transaction: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrTransactionNotFound
	}

	return nil
}

func (r *TransactionRepository) GetCurrentMonthIncome(ctx context.Context) (int64, error) {
	var income int64

	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE
			is_master_income = true
			AND date >= date_trunc('month', CURRENT_DATE)
			AND date < date_trunc('month', CURRENT_DATE) + interval '1 month'
	 `).Scan(&income)
	if err != nil {
		return 0, fmt.Errorf("get current month income: %w", err)
	}

	return income, nil
}
