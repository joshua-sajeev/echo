package allocations

import "errors"

var (
	ErrInvalidAmount                = errors.New("allocation amount must be greater than zero")
	ErrGoalNotFound                 = errors.New("goal not found")
	ErrInvalidGoalID                = errors.New("invalid goal id")
	ErrGoalArchived                 = errors.New("goal is archived")
	ErrNoGoalsConfigured            = errors.New("no goals configured")
	ErrInvalidAllocationPercentages = errors.New("goal allocation percentages must total 100")
)
