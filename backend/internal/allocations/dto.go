package allocations

// RunAllocationRequest for manual allocation to specific goal
type RunAllocationRequest struct {
	GoalID int64 `json:"goal_id" validate:"required,gt=0"`
	Amount int64 `json:"amount" validate:"required,min=1"`
}

// DistributeAllocationRequest for automatic allocation distributed by percentage
type DistributeAllocationRequest struct {
	Amount int64 `json:"amount" validate:"required,min=1"`
}
