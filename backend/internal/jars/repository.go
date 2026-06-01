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
