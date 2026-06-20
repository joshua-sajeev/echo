package allocations

import "errors"

var (
	ErrInvalidAmount                = errors.New("allocation amount must be greater than zero")
	ErrNoGoalsConfigured            = errors.New("no goals configured")
	ErrInvalidAllocationPercentages = errors.New("goal allocation percentages must total 100")
)
