// Package models
package models

import "time"

type TransactionAllocation struct {
	ID            int64 `json:"id"`
	TransactionID int64 `json:"transaction_id"`
	JarID         int64 `json:"jar_id"`
	Amount        int64 `json:"amount"`

	CreatedAt time.Time `json:"created_at"`
}
