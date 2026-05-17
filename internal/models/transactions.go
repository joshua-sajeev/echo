// Package models
package models

import (
	"time"
)

type Transaction struct {
	ID     int64     `json:"id"`
	Type   string    `json:"type"`
	Amount int64     `json:"amount"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`

	FromAccountID *int64 `json:"from_account_id"`
	ToAccountID   *int64 `json:"to_account_id"`

	Category    *string `json:"category"`
	TargetJarID *int64  `json:"target_jar_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
