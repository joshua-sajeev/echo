package goals

import "errors"

var (
	ErrGoalNameRequired          = errors.New("goal name is required")
	ErrInvalidGoalID             = errors.New("invalid goal id")
	ErrGoalNotFound              = errors.New("goal not found")
	ErrTargetAmountInvalid       = errors.New("target amount must be greater than 0")
	ErrProgressAmountInvalid     = errors.New("progress amount must be greater than 0")
	ErrGoalAlreadyCompleted      = errors.New("goal is already completed")
	ErrDeadlinePassed            = errors.New("deadline has already passed")
	ErrAllocationPercentageTotal = errors.New("total allocation percentage must equal 100")
	ErrInvalidAllocationCount    = errors.New("must have at least one new goal and existing goals")
	ErrGoalNotInRebalance        = errors.New("goal not found in rebalance request")
)
