package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshu-sajeev/echo/internal/models"
)

type JarRepository interface {
	List(ctx context.Context) ([]models.Jar, error)
}

type pgJarRepo struct {
	conn *pgx.Conn
}

func NewJarRepository(conn *pgx.Conn) JarRepository {
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
