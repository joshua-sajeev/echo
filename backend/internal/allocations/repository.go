package allocations

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	// Retrieve last month's leftover amount for the leisure jar
	GetLastMonthLeisureLeftover(ctx context.Context) (int64, error)

	// Check if automatic allocation was already executed this month
	IsAlreadyAllocatedThisMonth(ctx context.Context) (bool, error)
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

// GetLastMonthLeisureLeftover calculates the remaining balance/leftover for the leisure jar from last month
func (r *AllocationRepository) GetLastMonthLeisureLeftover(ctx context.Context) (int64, error) {
	// 1. Find the leisure jar
	var leisureID int64
	var leisureAllocType string
	var leisureValue int64
	err := r.conn.QueryRow(ctx, `
		SELECT id, allocation_type, value
		FROM jars
		WHERE LOWER(name) = 'leisure'
	`).Scan(&leisureID, &leisureAllocType, &leisureValue)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrLeisureJarNotFound
		}
		return 0, fmt.Errorf("get leisure jar: %w", err)
	}

	// 2. Fetch all jars to calculate remainder jar if needed
	rows, err := r.conn.Query(ctx, `SELECT id, allocation_type, value FROM jars`)
	if err != nil {
		return 0, fmt.Errorf("list jars for leftover: %w", err)
	}
	defer rows.Close()

	type jarInfo struct {
		id        int64
		allocType string
		value     int64
	}
	var allJars []jarInfo
	for rows.Next() {
		var ji jarInfo
		if err := rows.Scan(&ji.id, &ji.allocType, &ji.value); err != nil {
			return 0, fmt.Errorf("scan jar for leftover: %w", err)
		}
		allJars = append(allJars, ji)
	}

	// 3. Fetch last month's master income
	masterRows, err := r.conn.Query(ctx, `
		SELECT amount
		FROM transactions
		WHERE is_master_income = true
		  AND date >= date_trunc('month', CURRENT_DATE) - interval '1 month'
		  AND date < date_trunc('month', CURRENT_DATE)
	`)
	if err != nil {
		return 0, fmt.Errorf("query last month master income: %w", err)
	}
	defer masterRows.Close()

	var allocatedMasterIncome int64
	for masterRows.Next() {
		var amount int64
		if err := masterRows.Scan(&amount); err != nil {
			return 0, fmt.Errorf("scan last month master income: %w", err)
		}

		if leisureAllocType == "percentage" {
			allocatedMasterIncome += amount * leisureValue / 100
		} else if leisureAllocType == "remainder" {
			var percentageTotal int64
			for _, j := range allJars {
				if j.allocType == "percentage" {
					percentageTotal += amount * j.value / 100
				}
			}
			allocatedMasterIncome += amount - percentageTotal
		}
	}

	// 4. Fetch direct transactions for the leisure jar last month
	txRows, err := r.conn.Query(ctx, `
		SELECT type, amount
		FROM transactions
		WHERE jar_id = $1
		  AND is_master_income = false
		  AND type IN ('income', 'expense')
		  AND date >= date_trunc('month', CURRENT_DATE) - interval '1 month'
		  AND date < date_trunc('month', CURRENT_DATE)
	`, leisureID)
	if err != nil {
		return 0, fmt.Errorf("query last month direct transactions: %w", err)
	}
	defer txRows.Close()

	var directIncome int64
	var expenses int64

	for txRows.Next() {
		var txType string
		var amount int64
		if err := txRows.Scan(&txType, &amount); err != nil {
			return 0, fmt.Errorf("scan last month direct transaction: %w", err)
		}
		if txType == "income" {
			directIncome += amount
		} else if txType == "expense" {
			expenses += amount
		}
	}

	leftover := allocatedMasterIncome + directIncome - expenses
	if leftover < 0 {
		return 0, nil // Leftover cannot be negative for allocation purposes
	}
	return leftover, nil
}

// IsAlreadyAllocatedThisMonth checks if an automatic allocation transaction has run in the current month
func (r *AllocationRepository) IsAlreadyAllocatedThisMonth(ctx context.Context) (bool, error) {
	var exists bool
	err := r.conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 
			FROM goal_transactions 
			WHERE transaction_type = 'allocation' 
			  AND notes = 'Automatic allocation' 
			  AND created_at >= date_trunc('month', CURRENT_DATE)
		)
	`).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check if already allocated: %w", err)
	}
	return exists, nil
}

