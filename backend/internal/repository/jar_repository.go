// Package repository handles repository functions for the models
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type pgJarRepo struct {
	conn *pgxpool.Pool
}

func NewJarRepository(conn *pgxpool.Pool) *pgJarRepo {
	return &pgJarRepo{conn: conn}
}

func (r *pgJarRepo) List(ctx context.Context) ([]models.Jar, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT id, name, allocation_value, sort_order, is_system, created_at
		FROM jars
		ORDER BY sort_order ASC, id ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jars []models.Jar
	for rows.Next() {
		var j models.Jar
		if err := rows.Scan(&j.ID, &j.Name, &j.AllocationValue, &j.SortOrder, &j.IsSystem, &j.CreatedAt); err != nil {
			return nil, err
		}
		jars = append(jars, j)
	}
	return jars, rows.Err()
}

// ListWithSpend returns all jars plus a map of jar_id → this-month expenses.
func (r *pgJarRepo) ListWithSpend(ctx context.Context) ([]models.Jar, map[int64]float64, error) {
	jars, err := r.List(ctx)
	if err != nil {
		return nil, nil, err
	}

	rows, err := r.conn.Query(ctx, `
		SELECT jar_id, COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE type = 'expense'
		  AND jar_id IS NOT NULL
		  AND DATE_TRUNC('month', COALESCE(NULLIF(date, '0001-01-01'), created_at)) = DATE_TRUNC('month', NOW())
		GROUP BY jar_id
	`)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	spent := make(map[int64]float64)
	for rows.Next() {
		var id int64
		var amt float64
		if err := rows.Scan(&id, &amt); err != nil {
			return nil, nil, err
		}
		spent[id] = amt
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return jars, spent, nil
}

func (r *pgJarRepo) Update(ctx context.Context, id int64, name string, value float64) error {
	_, err := r.conn.Exec(ctx, `
		UPDATE jars SET name = $1, allocation_value = $2 WHERE id = $3
	`, name, value, id)
	return err
}

func (r *pgJarRepo) Create(ctx context.Context, name string, value float64, sortOrder int) (int64, error) {
	var id int64
	err := r.conn.QueryRow(ctx, `
		INSERT INTO jars (name, allocation_value, sort_order, is_system)
		VALUES ($1, $2, $3, false)
		RETURNING id
	`, name, value, sortOrder).Scan(&id)
	return id, err
}

func (r *pgJarRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.conn.Exec(ctx, `
		DELETE FROM jars WHERE id = $1 AND is_system = false
	`, id)
	return err
}

// EnsureDefaults creates the five system jars if they don't exist yet.
// Safe to call on every startup.
func (r *pgJarRepo) EnsureDefaults(ctx context.Context) error {
	defaults := []struct {
		name  string
		value float64
		order int
	}{
		{"Charity", 10, 1},
		{"SIP", 1000, 2},
		{"Chitty", 5000, 3},
		{"Necessities", 0, 4},
		{"Leisure", 0, 5},
	}

	for _, d := range defaults {
		_, err := r.conn.Exec(ctx, `
			INSERT INTO jars (name, allocation_value, sort_order, is_system)
			VALUES ($1, $2, $3, true)
			ON CONFLICT (name) DO NOTHING
		`, d.name, d.value, d.order)
		if err != nil {
			return err
		}
	}
	return nil
}
