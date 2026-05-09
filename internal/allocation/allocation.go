// Package allocation implements the jar allocation engine.
//
// Rules are hardcoded here — the DB stores only user-editable values
// (the numbers), not the allocation logic itself.
//
// Allocation order:
//  1. Charity       — percent_total:     round_down(total * value/100, 100)
//  2. SIP           — fixed:             value (₹1,000 default)
//  3. Chitty        — cap:               value is a spending cap, NOT deducted from pool
//  4. Necessities   — percent_remainder: round_up(pool - leisure_max, 500)
//  5. Leisure       — remainder:         pool - necessities  (≤ 10% of total)
//
// Any jar not matching a known name falls through as fixed allocation.
package allocation

import (
	"math"

	"github.com/joshu-sajeev/echo/internal/models"
)

// JarAllocation is a jar with its computed monthly figures.
type JarAllocation struct {
	models.Jar
	Allocated float64 // computed budget for this month
	Spent     float64 // actual expenses recorded against this jar
	Remaining float64 // Allocated - Spent (can be negative = overspent)
	IsCap     bool    // true for Chitty — Allocated is a cap, not a pre-allocation
}

// roundDown rounds x down to the nearest `unit`.
func roundDown(x float64, unit int) float64 {
	if unit <= 0 {
		return x
	}
	u := float64(unit)
	return math.Floor(x/u) * u
}

// roundUp rounds x up to the nearest `unit`.
func roundUp(x float64, unit int) float64 {
	if unit <= 0 {
		return x
	}
	u := float64(unit)
	return math.Ceil(x/u) * u
}

// Calculate computes the monthly allocation for each jar given the total
// master income for the month. spentByJar maps jar ID → amount spent.
//
// Jars are processed in sort_order. Unknown jar names are treated as fixed.
func Calculate(totalIncome float64, jars []models.Jar, spentByJar map[int64]float64) []JarAllocation {
	result := make([]JarAllocation, len(jars))

	// --- pass 1: resolve charity + fixed (SIP) allocations ----------------
	var poolDeductions float64 // total taken before necessities/leisure

	for i, j := range jars {
		result[i].Jar = j
		result[i].Spent = spentByJar[j.ID]

		switch j.Name {
		case "Charity":
			alloc := roundDown(totalIncome*j.AllocationValue/100, 100)
			result[i].Allocated = alloc
			poolDeductions += alloc

		case "SIP":
			result[i].Allocated = j.AllocationValue
			poolDeductions += j.AllocationValue

		case "Chitty":
			// cap jar — spending limit only, not deducted from pool
			result[i].Allocated = j.AllocationValue
			result[i].IsCap = true
			// deliberately NOT added to poolDeductions

		case "Necessities", "Leisure":
			// handled in pass 2
		default:
			// any user-added jar: treat as fixed, deduct from pool
			result[i].Allocated = j.AllocationValue
			poolDeductions += j.AllocationValue
		}
	}

	// --- pass 2: necessities + leisure from remainder pool -----------------
	pool := totalIncome - poolDeductions
	if pool < 0 {
		pool = 0
	}

	// leisure cap: floor of 10% of total income, to nearest ₹100
	leisureCap := roundDown(totalIncome*0.10, 100)
	if leisureCap > pool {
		leisureCap = pool
	}

	// necessities: pool minus leisure cap, rounded UP to nearest ₹500
	necessities := roundUp(pool-leisureCap, 500)
	if necessities > pool {
		necessities = pool
	}

	// actual leisure is what's genuinely left
	leisureActual := pool - necessities
	if leisureActual < 0 {
		leisureActual = 0
	}

	// write back
	for i, j := range jars {
		switch j.Name {
		case "Necessities":
			result[i].Allocated = necessities
		case "Leisure":
			result[i].Allocated = leisureActual
		}
		result[i].Remaining = result[i].Allocated - result[i].Spent
	}

	return result
}

// Summary holds the high-level monthly budget summary.
type Summary struct {
	TotalIncome    float64
	TotalAllocated float64 // sum of all non-cap allocations
	TotalSpent     float64
	Unallocated    float64 // should be 0 if rules are correct
}

// Summarise computes the summary from a slice of JarAllocation.
func Summarise(allocs []JarAllocation) Summary {
	var s Summary
	for _, a := range allocs {
		s.TotalSpent += a.Spent
		if !a.IsCap {
			s.TotalAllocated += a.Allocated
		}
	}
	return s
}
