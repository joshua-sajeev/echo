package models

import "time"

type Jar struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	TargetAmount float64   `json:"target_amount"`
	CreatedAt    time.Time `json:"created_at"`
}
