package goals

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoalRepository struct {
	conn *pgxpool.Pool
}

type GoalRepositoryInterface interface {
	Create(ctx context.Context, goal Goal) (int64, error)
	GetByID(ctx context.Context, id int64) (*Goal, error)
	List(ctx context.Context) ([]Goal, error)
	Update(ctx context.Context, goal Goal) error
	Archive(ctx context.Context, id int64) error
	AddProgress(ctx context.Context, id int64, amount int64) error
	Exists(ctx context.Context, id int64) (bool, error)
	Restore(ctx context.Context, id int64) error
}

var _ GoalRepositoryInterface = (*GoalRepository)(nil)

// NewGoalRepository creates a new goal repository
func NewGoalRepository(conn *pgxpool.Pool) *GoalRepository {
	return &GoalRepository{conn: conn}
}

// Create inserts a new goal and returns its ID
func (r *GoalRepository) Create(ctx context.Context, goal Goal) (int64, error) {
	var id int64

	err := r.conn.QueryRow(
		ctx,
		`INSERT INTO goals (
		name,
		target_amount,
		saved_amount,
		deadline,
		allocation_percentage
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`,
		goal.Name,
		goal.TargetAmount,
		goal.SavedAmount,
		goal.Deadline,
		goal.AllocationPercentage,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create goal: %w", err)
	}

	return id, nil
}

// GetByID retrieves a goal by its ID
func (r *GoalRepository) GetByID(ctx context.Context, id int64) (*Goal, error) {
	var goal Goal

	err := r.conn.QueryRow(
		ctx,
		`
	SELECT id,
		name,
		target_amount,
		saved_amount,
		deadline,
		allocation_percentage,
		is_archived,
		created_at,
		updated_at
	FROM goals
	WHERE id = $1
	`,
		id,
	).Scan(
		&goal.ID,
		&goal.Name,
		&goal.TargetAmount,
		&goal.SavedAmount,
		&goal.Deadline,
		&goal.AllocationPercentage,
		&goal.IsArchived,
		&goal.CreatedAt,
		&goal.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrGoalNotFound
		}
		return nil, fmt.Errorf("get goal by id: %w", err)
	}

	return &goal, nil
}

// List returns all goals ordered by creation date
func (r *GoalRepository) List(ctx context.Context) ([]Goal, error) {
	rows, err := r.conn.Query(
		ctx,
		`
	SELECT
		id,
		name,
		target_amount,
		saved_amount,
		deadline,
		allocation_percentage,
		is_archived,
		created_at,
		updated_at
	FROM goals
	WHERE is_archived = FALSE
	ORDER BY created_at DESC
	`,
	)
	if err != nil {
		return nil, fmt.Errorf("list goals: %w", err)
	}
	defer rows.Close()

	goals := make([]Goal, 0)
	for rows.Next() {
		var goal Goal
		if err := rows.Scan(
			&goal.ID,
			&goal.Name,
			&goal.TargetAmount,
			&goal.SavedAmount,
			&goal.Deadline,
			&goal.AllocationPercentage,
			&goal.IsArchived,
			&goal.CreatedAt,
			&goal.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan goal: %w", err)
		}
		goals = append(goals, goal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating goals: %w", err)
	}

	return goals, nil
}

// Update modifies an existing goal
func (r *GoalRepository) Update(ctx context.Context, goal Goal) error {
	tag, err := r.conn.Exec(
		ctx,
		`UPDATE goals
		 SET name = $1, target_amount = $2, deadline = $3, allocation_percentage = $4, updated_at = NOW()
		 WHERE id = $5`,
		goal.Name,
		goal.TargetAmount,
		goal.Deadline,
		goal.AllocationPercentage,
		goal.ID,
	)
	if err != nil {
		return fmt.Errorf("update goal: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGoalNotFound
	}

	return nil
}

// Archives a goal
func (r *GoalRepository) Archive(ctx context.Context, id int64) error {
	tag, err := r.conn.Exec(
		ctx,
		`
		UPDATE goals
		SET
			is_archived = TRUE,
			updated_at = NOW()
		WHERE id = $1
		  AND is_archived = FALSE
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("archive goal: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGoalNotFound
	}

	return nil
}

// AddProgress adds to the saved_amount and updates updated_at timestamp, creating a transaction record
func (r *GoalRepository) AddProgress(ctx context.Context, id int64, amount int64) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(
		ctx,
		`UPDATE goals
		 SET saved_amount = saved_amount + $1, updated_at = NOW()
		 WHERE id = $2`,
		amount,
		id,
	)
	if err != nil {
		return fmt.Errorf("add progress update: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGoalNotFound
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO goal_transactions (
			goal_id,
			amount,
			transaction_type,
			notes
		)
		VALUES ($1, $2, 'manual_contribution', 'Manual contribution')`,
		id,
		amount,
	)
	if err != nil {
		return fmt.Errorf("insert goal transaction for progress: %w", err)
	}

	return tx.Commit(ctx)
}

// Exists checks if a goal exists
func (r *GoalRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool

	err := r.conn.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM goals WHERE id = $1)`,
		id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check goal exists: %w", err)
	}

	return exists, nil
}

// Restore restores a archived goal
func (r *GoalRepository) Restore(ctx context.Context, id int64) error {
	tag, err := r.conn.Exec(
		ctx,
		`
		UPDATE goals
		SET
			is_archived = FALSE,
			updated_at = NOW()
		WHERE id = $1
		  AND is_archived = TRUE
		`,
		id,
	)
	if err != nil {
		return fmt.Errorf("restore goal: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return ErrGoalNotFound
	}

	return nil
}
