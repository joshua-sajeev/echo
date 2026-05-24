// Package models contains details about different tables
package models

import "time"

type Account struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	IsArchived bool   `json:"is_archived"`

	CreatedAt time.Time `json:"created_at"`
}

type AccountWithBalance struct {
	Account
	Balance float64
}
