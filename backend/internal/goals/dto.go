package goals

import "time"

type CreateGoalRequest struct {
	Name                 string     `json:"name" validate:"required,min=1,max=100"`
	TargetAmount         int64      `json:"target_amount" validate:"required,gt=0"`
	AllocationPercentage int64      `json:"allocation_percentage" validate:"required,gte=0,lte=100"`
	Deadline             *time.Time `json:"deadline" validate:"omitempty"`
}

type UpdateGoalRequest struct {
	Name                 *string    `json:"name" validate:"omitempty,min=1,max=100"`
	TargetAmount         *int64     `json:"target_amount" validate:"omitempty,gt=0"`
	AllocationPercentage *int64     `json:"allocation_percentage" validate:"omitempty,gte=0,lte=100"`
	Deadline             *time.Time `json:"deadline" validate:"omitempty"`
}

type AddProgressRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}

// GoalAllocationItem represents a goal with its new allocation percentage
type GoalAllocationItem struct {
	ID                   *int64     `json:"id" validate:"omitempty"` // nil for new goal
	Name                 string     `json:"name" validate:"required,min=1,max=100"`
	TargetAmount         int64      `json:"target_amount" validate:"required,gt=0"`
	AllocationPercentage int64      `json:"allocation_percentage" validate:"required,gte=0,lte=100"`
	Deadline             *time.Time `json:"deadline" validate:"omitempty"`
}

// CreateGoalWithRebalanceRequest creates a new goal and rebalances all others
type CreateGoalWithRebalanceRequest struct {
	Goals []GoalAllocationItem `json:"goals" validate:"required"` // At least new goal + 1 existing
}

// AllocationChange shows what changed for a goal
type AllocationChange struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	OldPercentage int64  `json:"old_percentage"`
	NewPercentage int64  `json:"new_percentage"`
	IsNew         bool   `json:"is_new"`
}

// RebalanceSummary shows the allocation changes made
type RebalanceSummary struct {
	Changes []AllocationChange `json:"changes"`
	Message string             `json:"message"`
}
