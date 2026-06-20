package allocations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/goals"
)

type AllocationRepositoryInterface interface {
	// Manual: allocate to one goal
	RunManual(
		ctx context.Context,
		goalID int64,
		amount int64,
	) error

	// Automatic: distribute across all goals by percentage
	DistributeAutomatic(
		ctx context.Context,
		amount int64,
		goalsList []goals.Goal,
	) error
}

type AllocationRepository struct {
	conn *pgxpool.Pool
}

func NewAllocationRepository(
	conn *pgxpool.Pool,
) *AllocationRepository {
	return &AllocationRepository{
		conn: conn,
	}
}

// RunManual allocates amount to a specific goal
func (r *AllocationRepository) RunManual(
	ctx context.Context,
	goalID int64,
	amount int64,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update goal saved_amount
	_, err = tx.Exec(
		ctx,
		`UPDATE goals
		 SET saved_amount = saved_amount + $1,
		     updated_at = NOW()
		 WHERE id = $2`,
		amount,
		goalID,
	)
	if err != nil {
		return fmt.Errorf("update goal: %w", err)
	}

	// Create goal_transaction record
	_, err = tx.Exec(
		ctx,
		`INSERT INTO goal_transactions (
			goal_id,
			amount,
			transaction_type,
			notes
		)
		VALUES ($1, $2, 'allocation', 'Manual allocation')`,
		goalID,
		amount,
	)
	if err != nil {
		return fmt.Errorf("insert goal transaction: %w", err)
	}

	return tx.Commit(ctx)
}

// DistributeAutomatic distributes amount across all goals by their allocation_percentage
func (r *AllocationRepository) DistributeAutomatic(
	ctx context.Context,
	amount int64,
	goalsList []goals.Goal,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var allocatedTotal int64

	for i, goal := range goalsList {
		var allocation int64

		// Last goal gets remainder (handles rounding)
		if i == len(goalsList)-1 {
			allocation = amount - allocatedTotal
		} else {
			allocation = (amount * goal.AllocationPercentage) / 100
			allocatedTotal += allocation
		}

		// Update goal saved_amount
		_, err = tx.Exec(
			ctx,
			`UPDATE goals
			 SET saved_amount = saved_amount + $1,
			     updated_at = NOW()
			 WHERE id = $2`,
			allocation,
			goal.ID,
		)
		if err != nil {
			return fmt.Errorf("update goal: %w", err)
		}

		// Create goal_transaction record
		_, err = tx.Exec(
			ctx,
			`INSERT INTO goal_transactions (
				goal_id,
				amount,
				transaction_type,
				notes
			)
			VALUES ($1, $2, 'allocation', 'Automatic allocation')`,
			goal.ID,
			allocation,
		)
		if err != nil {
			return fmt.Errorf("insert goal transaction: %w", err)
		}
	}

	return tx.Commit(ctx)
}
