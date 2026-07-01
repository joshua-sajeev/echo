package allocations

// RunAllocationRequest for manual allocation to specific goal
type RunAllocationRequest struct {
	GoalID int64 `json:"goal_id" validate:"required,gt=0"`
	Amount int64 `json:"amount" validate:"required,min=1"`
}

// DistributeAllocationRequest for automatic allocation distributed by percentage
type DistributeAllocationRequest struct {
	Type   string `json:"type" validate:"required"` // "automatic_splitting" or "leisure_leftover"
	Amount *int64 `json:"amount" validate:"omitempty,min=1"`
}
