package models

import "time"

type TxTemplate struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	JarID           *int64    `json:"jar_id"`
	JarName         string    `json:"jar_name"`
	Amount          float64   `json:"amount"`
	Type            string    `json:"type"` // expense | income | transfer
	FromAccountID   *int64    `json:"from_account_id"`
	FromAccountName string    `json:"from_account_name"`
	ToAccountID     *int64    `json:"to_account_id"`
	ToAccountName   string    `json:"to_account_name"`
	IsMasterIncome  bool      `json:"is_master_income"`
	CreatedAt       time.Time `json:"created_at"`
}
