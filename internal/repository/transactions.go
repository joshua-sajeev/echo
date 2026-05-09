package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshu-sajeev/echo/internal/models"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx models.Transaction) error
	List(ctx context.Context) ([]models.Transaction, error)
}

type pgTransactionRepo struct {
	conn *pgx.Conn
}

func NewTransactionRepository(conn *pgx.Conn) TransactionRepository {
	return &pgTransactionRepo{conn: conn}
}

func (r *pgTransactionRepo) Create(ctx context.Context, tx models.Transaction) error {
	_, err := r.conn.Exec(ctx, `
		INSERT INTO transactions
			(type, amount, name, date, from_account_id, to_account_id, jar_id, is_master_income)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
	`, tx.Type, tx.Amount, tx.Name, tx.Date,
		tx.FromAccountID, tx.ToAccountID, tx.JarID, tx.IsMasterIncome)
	return err
}

func (r *pgTransactionRepo) List(ctx context.Context) ([]models.Transaction, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, name, amount, type, date, created_at,
		       from_account_id, to_account_id, jar_id, is_master_income
		FROM transactions
		ORDER BY date DESC, id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		if err := rows.Scan(
			&tx.ID, &tx.Name, &tx.Amount, &tx.Type, &tx.Date, &tx.CreatedAt,
			&tx.FromAccountID, &tx.ToAccountID, &tx.JarID, &tx.IsMasterIncome,
		); err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, rows.Err()
}
