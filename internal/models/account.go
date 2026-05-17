// Package models
package models

import "time"

type Account struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	ArchivedAt *time.Time `json:"archived_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type AccountWithBalance struct {
	Account
	Balance int64
}
