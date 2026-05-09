package models

import "time"

type Transaction struct {
	ID   int64  `json:"id"`
	Type string `json:"type"` // income | expense | transfer

	Amount float64   `json:"amount"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`

	FromAccountID *int64 `json:"from_account_id"`
	ToAccountID   *int64 `json:"to_account_id"`

	Category    *string `json:"category"`
	SubCategory *string `json:"sub_category"`

	JarID *int64 `json:"jar_id"`

	IsMasterIncome bool `json:"is_master_income"`

	CreatedAt time.Time `json:"created_at"`
}
