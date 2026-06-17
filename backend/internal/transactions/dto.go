package transactions

import "time"

type CreateTransactionRequest struct {
	Name           string    `json:"name"            validate:"required,min=1,max=100"`
	Type           string    `json:"type"            validate:"required,oneof=income expense transfer"`
	Amount         int64     `json:"amount"          validate:"required,gt=0"`
	Date           time.Time `json:"date"            validate:"required"`
	FromAccountID  *int64    `json:"from_account_id" validate:"omitempty,gt=0"`
	ToAccountID    *int64    `json:"to_account_id"   validate:"omitempty,gt=0"`
	Category       *string   `json:"category"        validate:"omitempty,min=1,max=50"`
	JarID          *int64    `json:"jar_id"          validate:"omitempty,gt=0"`
	IsMasterIncome bool      `json:"is_master_income"`
}

type UpdateTransactionRequest struct {
	Name           *string    `json:"name"            validate:"omitempty,min=1,max=100"`
	Type           *string    `json:"type"            validate:"omitempty,oneof=income expense transfer"`
	Amount         *int64     `json:"amount"          validate:"omitempty,gt=0"`
	Date           *time.Time `json:"date"            validate:"omitempty"`
	FromAccountID  *int64     `json:"from_account_id" validate:"omitempty,gt=0"`
	ToAccountID    *int64     `json:"to_account_id"   validate:"omitempty,gt=0"`
	Category       *string    `json:"category"        validate:"omitempty,min=1,max=50"`
	JarID          *int64     `json:"jar_id"          validate:"omitempty,gt=0"`
	IsMasterIncome *bool      `json:"is_master_income" validate:"omitempty"`
}

type TransactionListItem struct {
	ID     int64     `json:"id"`
	Type   string    `json:"type"`
	Amount int64     `json:"amount"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`

	FromAccountID   *int64  `json:"from_account_id"`
	FromAccountName *string `json:"from_account_name"`

	ToAccountID   *int64  `json:"to_account_id"`
	ToAccountName *string `json:"to_account_name"`

	JarID   *int64  `json:"jar_id"`
	JarName *string `json:"jar_name"`

	Category       *string   `json:"category"`
	IsMasterIncome bool      `json:"is_master_income"`
	CreatedAt      time.Time `json:"created_at"`
}
