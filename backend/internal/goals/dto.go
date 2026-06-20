package goals

import "time"

type CreateGoalRequest struct {
	Name         string     `json:"name" validate:"required,min=1,max=100"`
	TargetAmount int64      `json:"target_amount" validate:"required,gt=0"`
	Deadline     *time.Time `json:"deadline" validate:"omitempty"`
}

type UpdateGoalRequest struct {
	Name         *string    `json:"name" validate:"omitempty,min=1,max=100"`
	TargetAmount *int64     `json:"target_amount" validate:"omitempty,gt=0"`
	Deadline     *time.Time `json:"deadline" validate:"omitempty"`
}

type AddProgressRequest struct {
	Amount int64 `json:"amount" validate:"required,gt=0"`
}
