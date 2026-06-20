package allocations

import "time"

type GoalTransaction struct {
	ID              int64     `json:"id"`
	GoalID          int64     `json:"goal_id"`
	Amount          int64     `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}
