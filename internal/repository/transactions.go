package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type TransactionRow struct {
	models.Transaction
	FromAccountName string
	ToAccountName   string
	JarName         string
}

type TransactionRepository interface {
	Create(ctx context.Context, tx models.Transaction) error
	List(ctx context.Context) ([]TransactionRow, error)
}

type pgTransactionRepo struct {
	conn *pgxpool.Pool
}

func NewTransactionRepository(conn *pgxpool.Pool) TransactionRepository {
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

func (r *pgTransactionRepo) List(ctx context.Context) ([]TransactionRow, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT t.id, t.name, t.amount, t.type, t.date, t.created_at,
		       t.from_account_id, t.to_account_id, t.jar_id, t.is_master_income,
		       COALESCE(fa.name, ''), COALESCE(ta.name, ''),
		       COALESCE(j.name, '')
		FROM transactions t
		LEFT JOIN accounts fa ON fa.id = t.from_account_id
		LEFT JOIN accounts ta ON ta.id = t.to_account_id
		LEFT JOIN jars j ON j.id = t.jar_id
		ORDER BY t.date DESC, t.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []TransactionRow
	for rows.Next() {
		var row TransactionRow
		if err := rows.Scan(
			&row.ID, &row.Name, &row.Amount, &row.Type, &row.Date, &row.CreatedAt,
			&row.FromAccountID, &row.ToAccountID, &row.JarID, &row.IsMasterIncome,
			&row.FromAccountName, &row.ToAccountName, &row.JarName,
		); err != nil {
			return nil, err
		}
		txs = append(txs, row)
	}
	return txs, rows.Err()
}
