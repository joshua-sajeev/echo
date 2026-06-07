// Package jars handles jars related code
package jars

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JarRepository struct {
	conn *pgxpool.Pool
}

type JarRepositoryInterface interface {
	Create(ctx context.Context, jar Jar) (int64, error)
	List(ctx context.Context) ([]Jar, error)
	GetByID(ctx context.Context, id int64) (Jar, error)
	Update(ctx context.Context, jar Jar) error
	Delete(ctx context.Context, id int64) error

	// GetAllJarBalances returns a map of jar ID → running balance, calculated
	// on the fly from all transactions.
	GetAllJarBalances(ctx context.Context) (map[int64]int64, error)

	// GetSpentThisMonthPerJar returns a map of jar ID → total expenses
	// charged to that jar in the current calendar month.
	GetSpentThisMonthPerJar(ctx context.Context) (map[int64]int64, error)
}

var _ JarRepositoryInterface = (*JarRepository)(nil)

func NewJarRepository(conn *pgxpool.Pool) *JarRepository {
	return &JarRepository{conn: conn}
}

func (r *JarRepository) Create(ctx context.Context, jar Jar) (int64, error) {
	var id int64

	err := r.conn.QueryRow(
		ctx,
		`
		INSERT INTO jars (name, allocation_type, value)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		jar.Name,
		jar.AllocationType,
		jar.Value,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create jar: %w", err)
	}

	return id, nil
}

func (r *JarRepository) List(ctx context.Context) ([]Jar, error) {
	rows, err := r.conn.Query(
		ctx,
		`
		SELECT id, name, allocation_type, value, created_at
		FROM jars
		ORDER BY id ASC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list jars: %w", err)
	}
	defer rows.Close()

	jars := make([]Jar, 0)

	for rows.Next() {
		var jar Jar

		err := rows.Scan(
			&jar.ID,
			&jar.Name,
			&jar.AllocationType,
			&jar.Value,
			&jar.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan jar: %w", err)
		}

		jars = append(jars, jar)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate jars: %w", err)
	}

	return jars, nil
}

func (r *JarRepository) GetByID(ctx context.Context, id int64) (Jar, error) {
	var jar Jar

	err := r.conn.QueryRow(
		ctx,
		`
		SELECT id, name, allocation_type, value, created_at
		FROM jars
		WHERE id = $1
		`,
		id,
	).Scan(
		&jar.ID,
		&jar.Name,
		&jar.AllocationType,
		&jar.Value,
		&jar.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Jar{}, ErrJarNotFound
		}

		return Jar{}, fmt.Errorf("get jar by id: %w", err)
	}

	return jar, nil
}

func (r *JarRepository) Update(ctx context.Context, jar Jar) error {
	tag, err := r.conn.Exec(
		ctx,
		`
		UPDATE jars
		SET
			name = $1,
			allocation_type = $2,
			value = $3
		WHERE id = $4
		`,
		jar.Name,
		jar.AllocationType,
		jar.Value,
		jar.ID,
	)
	if err != nil {
		return fmt.Errorf("update jar: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrJarNotFound
	}

	return nil
}

func (r *JarRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.conn.Exec(ctx, `DELETE FROM jars WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete jar: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrJarNotFound
	}

	return nil
}

// GetAllJarBalances calculates the running balance for every jar in a single
// query. It does three passes in SQL and merges them in Go.
func (r *JarRepository) GetAllJarBalances(ctx context.Context) (map[int64]int64, error) {
	// ── Step 1: load all jars so we know allocation types / values ──────────
	jars, err := r.List(ctx)
	if err != nil {
		return nil, err
	}
	if len(jars) == 0 {
		return map[int64]int64{}, nil
	}

	balances := make(map[int64]int64, len(jars))
	for _, j := range jars {
		balances[j.ID] = 0
	}

	// ── Step 2: master-income allocations ───────────────────────────────────
	masterRows, err := r.conn.Query(ctx, `
		SELECT amount
		FROM transactions
		WHERE is_master_income = true
	`)
	if err != nil {
		return nil, fmt.Errorf("query master income: %w", err)
	}
	defer masterRows.Close()

	var remainderJarID int64 = -1
	for _, j := range jars {
		if j.AllocationType == AllocationRemainder {
			remainderJarID = j.ID
		}
	}

	for masterRows.Next() {
		var amount int64
		if err := masterRows.Scan(&amount); err != nil {
			return nil, fmt.Errorf("scan master income: %w", err)
		}

		var percentageTotal int64
		for _, j := range jars {
			if j.AllocationType == AllocationPercentage {
				share := amount * j.Value / 100
				balances[j.ID] += share
				percentageTotal += share
			}
		}
		if remainderJarID > 0 {
			balances[remainderJarID] += amount - percentageTotal
		}
	}
	if err := masterRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate master income: %w", err)
	}

	// ── Step 3: direct income & expenses tagged to a jar ────────────────────
	txRows, err := r.conn.Query(ctx, `
		SELECT jar_id, type, amount
		FROM transactions
		WHERE jar_id IS NOT NULL
		  AND is_master_income = false
		  AND type IN ('income', 'expense')
	`)
	if err != nil {
		return nil, fmt.Errorf("query jar transactions: %w", err)
	}
	defer txRows.Close()

	for txRows.Next() {
		var jarID int64
		var txType string
		var amount int64

		if err := txRows.Scan(&jarID, &txType, &amount); err != nil {
			return nil, fmt.Errorf("scan jar transaction: %w", err)
		}

		if _, ok := balances[jarID]; !ok {
			continue
		}

		switch txType {
		case "income":
			balances[jarID] += amount
		case "expense":
			balances[jarID] -= amount
		}
	}
	if err := txRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate jar transactions: %w", err)
	}

	return balances, nil
}

// GetSpentThisMonthPerJar returns a map of jar ID → total expenses
// charged to that jar in the current calendar month.
func (r *JarRepository) GetSpentThisMonthPerJar(ctx context.Context) (map[int64]int64, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT jar_id, COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE
			type = 'expense'
			AND jar_id IS NOT NULL
			AND date >= date_trunc('month', CURRENT_DATE)
			AND date < date_trunc('month', CURRENT_DATE) + interval '1 month'
		GROUP BY jar_id
	`)
	if err != nil {
		return nil, fmt.Errorf("get spent this month per jar: %w", err)
	}
	defer rows.Close()

	result := make(map[int64]int64)
	for rows.Next() {
		var jarID, spent int64
		if err := rows.Scan(&jarID, &spent); err != nil {
			return nil, fmt.Errorf("scan spent this month: %w", err)
		}
		result[jarID] = spent
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate spent this month: %w", err)
	}
	return result, nil
}
