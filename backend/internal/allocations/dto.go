package allocations

type RunAllocationRequest struct {
	Amount int64 `json:"amount" validate:"required,min=1"`
}
