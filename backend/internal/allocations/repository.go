package allocations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshu-sajeev/echo/internal/goals"
)

type AllocationRepositoryInterface interface {
	RunAllocation(
		ctx context.Context,
		amount int64,
		goals []goals.Goal,
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

func (r *AllocationRepository) RunAllocation(
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

		if i == len(goalsList)-1 {
			allocation = amount - allocatedTotal
		} else {
			allocation = (amount * goal.AllocationPercentage) / 100

			allocatedTotal += allocation
		}

		_, err = tx.Exec(
			ctx,
			`
			UPDATE goals
			SET
				saved_amount = saved_amount + $1,
				updated_at = NOW()
			WHERE id = $2
			`,
			allocation,
			goal.ID,
		)
		if err != nil {
			return fmt.Errorf("update goal: %w", err)
		}

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO goal_transactions (
				goal_id,
				amount,
				transaction_type,
				notes
			)
			VALUES (
				$1,
				$2,
				'allocation',
				'Automatic leisure allocation'
			)
			`,
			goal.ID,
			allocation,
		)
		if err != nil {
			return fmt.Errorf(
				"insert goal transaction: %w",
				err,
			)
		}
	}

	return tx.Commit(ctx)
}
