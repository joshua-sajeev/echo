package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/models"
)

type JarWithSpend struct {
	models.Jar
	MonthlySpend float64
}

type JarRepository interface {
	List(ctx context.Context) ([]models.Jar, error)
	ListWithMonthlySpend(ctx context.Context) ([]JarWithSpend, error)
}

type pgJarRepo struct {
	conn *pgxpool.Pool
}

func NewJarRepository(conn *pgxpool.Pool) JarRepository {
	return &pgJarRepo{conn: conn}
}

func (r *pgJarRepo) List(ctx context.Context) ([]models.Jar, error) {
	rows, err := r.conn.Query(ctx,
		`SELECT id, name, target_amount, created_at FROM jars ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jars []models.Jar
	for rows.Next() {
		var j models.Jar
		if err := rows.Scan(&j.ID, &j.Name, &j.TargetAmount, &j.CreatedAt); err != nil {
			return nil, err
		}
		jars = append(jars, j)
	}
	return jars, rows.Err()
}

func (r *pgJarRepo) ListWithMonthlySpend(ctx context.Context) ([]JarWithSpend, error) {
	rows, err := r.conn.Query(ctx, `
        SELECT
            j.id, j.name, j.target_amount, j.created_at,
            COALESCE(SUM(CASE
                WHEN t.type = 'expense'
                AND DATE_TRUNC('month', COALESCE(NULLIF(t.date, '0001-01-01'), t.created_at)) = DATE_TRUNC('month', NOW())
                THEN t.amount ELSE 0 END), 0) AS monthly_spend
        FROM jars j
        LEFT JOIN transactions t ON t.jar_id = j.id
        GROUP BY j.id, j.name, j.target_amount, j.created_at
        ORDER BY j.id ASC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jars []JarWithSpend
	for rows.Next() {
		var j JarWithSpend
		if err := rows.Scan(&j.ID, &j.Name, &j.TargetAmount, &j.CreatedAt, &j.MonthlySpend); err != nil {
			return nil, err
		}
		jars = append(jars, j)
	}
	return jars, rows.Err()
}
