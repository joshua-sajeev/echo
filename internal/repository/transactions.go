package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type TransactionRow struct {
	models.Transaction
	FromAccountName string
	ToAccountName   string
	JarName         string
}

type Stats struct {
	TotalBalance        float64
	MonthlyIncome       float64
	MonthlyExpenses     float64
	Savings             float64
	MonthlyMasterIncome float64
}

// TxFilters holds all optional filters for ListAll.
type TxFilters struct {
	Type      string  // income | expense | transfer | ""
	Search    string  // name ILIKE
	AccountID int64   // from or to account id (0 = any)
	JarID     int64   // jar id (0 = any)
	AmountMin float64 // 0 = no lower bound
	AmountMax float64 // 0 = no upper bound
	DateFrom  string  // YYYY-MM-DD or ""
	DateTo    string  // YYYY-MM-DD or ""
}

type TransactionRepository interface {
	Create(ctx context.Context, tx models.Transaction) error
	List(ctx context.Context) ([]TransactionRow, error)
	ListAll(ctx context.Context, f TxFilters) ([]TransactionRow, error)
	Get(ctx context.Context, id int64) (TransactionRow, error)
	Update(ctx context.Context, tx models.Transaction) error
	Delete(ctx context.Context, id int64) error
	Stats(ctx context.Context) (Stats, error)
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

const txSelectCols = `
	SELECT t.id, t.name, t.amount, t.type, t.date, t.created_at,
	       t.from_account_id, t.to_account_id, t.jar_id, t.is_master_income,
	       COALESCE(fa.name, ''), COALESCE(ta.name, ''),
	       COALESCE(j.name, '')
	FROM transactions t
	LEFT JOIN accounts fa ON fa.id = t.from_account_id
	LEFT JOIN accounts ta ON ta.id = t.to_account_id
	LEFT JOIN jars j ON j.id = t.jar_id
`

func scanTxRow(row interface{ Scan(...any) error }) (TransactionRow, error) {
	var r TransactionRow
	err := row.Scan(
		&r.ID, &r.Name, &r.Amount, &r.Type, &r.Date, &r.CreatedAt,
		&r.FromAccountID, &r.ToAccountID, &r.JarID, &r.IsMasterIncome,
		&r.FromAccountName, &r.ToAccountName, &r.JarName,
	)
	return r, err
}

func (r *pgTransactionRepo) List(ctx context.Context) ([]TransactionRow, error) {
	rows, err := r.conn.Query(ctx, txSelectCols+`
		ORDER BY t.date DESC, t.id DESC
		LIMIT 7
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []TransactionRow
	for rows.Next() {
		row, err := scanTxRow(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, row)
	}
	return txs, rows.Err()
}

func (r *pgTransactionRepo) ListAll(ctx context.Context, f TxFilters) ([]TransactionRow, error) {
	query := txSelectCols + ` WHERE 1=1`
	args := []any{}

	if f.Type != "" && f.Type != "all" {
		args = append(args, f.Type)
		query += fmt.Sprintf(` AND t.type = $%d`, len(args))
	}
	if f.Search != "" {
		args = append(args, "%"+f.Search+"%")
		query += fmt.Sprintf(` AND t.name ILIKE $%d`, len(args))
	}
	if f.AccountID > 0 {
		args = append(args, f.AccountID)
		query += fmt.Sprintf(` AND (t.from_account_id = $%d OR t.to_account_id = $%d)`, len(args), len(args))
	}
	if f.JarID > 0 {
		args = append(args, f.JarID)
		query += fmt.Sprintf(` AND t.jar_id = $%d`, len(args))
	}
	if f.AmountMin > 0 {
		args = append(args, f.AmountMin)
		query += fmt.Sprintf(` AND t.amount >= $%d`, len(args))
	}
	if f.AmountMax > 0 {
		args = append(args, f.AmountMax)
		query += fmt.Sprintf(` AND t.amount <= $%d`, len(args))
	}
	if f.DateFrom != "" {
		args = append(args, f.DateFrom)
		query += fmt.Sprintf(` AND t.date >= $%d`, len(args))
	}
	if f.DateTo != "" {
		args = append(args, f.DateTo)
		query += fmt.Sprintf(` AND t.date <= $%d`, len(args))
	}
	query += ` ORDER BY t.date DESC, t.id DESC`

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []TransactionRow
	for rows.Next() {
		row, err := scanTxRow(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, row)
	}
	return txs, rows.Err()
}

func (r *pgTransactionRepo) Get(ctx context.Context, id int64) (TransactionRow, error) {
	row := r.conn.QueryRow(ctx, txSelectCols+` WHERE t.id = $1`, id)
	return scanTxRow(row)
}

func (r *pgTransactionRepo) Update(ctx context.Context, tx models.Transaction) error {
	_, err := r.conn.Exec(ctx, `
        UPDATE transactions
        SET name=$1, amount=$2, date=$3,
            from_account_id=$4, to_account_id=$5,
            jar_id=$6, is_master_income=$7
        WHERE id=$8
    `, tx.Name, tx.Amount, tx.Date,
		tx.FromAccountID, tx.ToAccountID,
		tx.JarID, tx.IsMasterIncome, tx.ID)
	return err
}

func (r *pgTransactionRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx, `DELETE FROM transactions WHERE id=$1`, id)
	return err
}

func (r *pgTransactionRepo) Stats(ctx context.Context) (Stats, error) {
	var s Stats

	err := r.conn.QueryRow(ctx, `
        SELECT
            COALESCE(SUM(CASE WHEN type = 'income'  THEN amount ELSE 0 END), 0) -
            COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0)
                AS total_balance,

            COALESCE(SUM(CASE
                WHEN type = 'income'
                AND DATE_TRUNC('month', COALESCE(NULLIF(date, '0001-01-01'), created_at)) = DATE_TRUNC('month', NOW())
                THEN amount ELSE 0 END), 0) AS monthly_income,

            COALESCE(SUM(CASE
                WHEN type = 'expense'
                AND DATE_TRUNC('month', COALESCE(NULLIF(date, '0001-01-01'), created_at)) = DATE_TRUNC('month', NOW())
                THEN amount ELSE 0 END), 0) AS monthly_expenses,

            COALESCE(SUM(CASE
                WHEN type = 'income'
                AND is_master_income = true
                AND DATE_TRUNC('month', COALESCE(NULLIF(date, '0001-01-01'), created_at)) = DATE_TRUNC('month', NOW())
                THEN amount ELSE 0 END), 0) AS monthly_master_income

        FROM transactions
    `).Scan(&s.TotalBalance, &s.MonthlyIncome, &s.MonthlyExpenses, &s.MonthlyMasterIncome)
	if err != nil {
		return s, err
	}

	s.Savings = s.MonthlyIncome - s.MonthlyExpenses
	return s, nil
}
