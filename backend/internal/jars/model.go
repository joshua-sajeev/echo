// Package jars
package jars

import "time"

type AllocationType string

const (
	AllocationPercentage AllocationType = "percentage"
	AllocationRemainder  AllocationType = "remainder"
)

type Jar struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	AllocationType AllocationType `json:"allocation_type"`

	// Used only when AllocationType == percentage
	Value int64 `json:"value"`

	CreatedAt time.Time `json:"created_at"`
}

type JarWithAllocation struct {
	Jar
	// AllocatedAmount is the share of the current month's master income for this jar.
	// For percentage jars: income * value / 100
	// For remainder jars: income - sum of all percentage allocations
	AllocatedAmount int64 `json:"allocated_amount"`

	// Balance is the running total across all time:
	//   + all master income allocations to this jar
	//   + all direct income transactions tagged to this jar
	//   - all expense transactions tagged to this jar
	Balance int64 `json:"balance"`

	// SpentThisMonth is the total expenses charged to this jar in the current
	// calendar month. Used to calculate progress: how much of this month's
	// allocation has been consumed.
	SpentThisMonth int64 `json:"spent_this_month"`
}
