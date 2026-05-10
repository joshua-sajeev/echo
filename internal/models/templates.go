package models

import "time"

type TxTemplate struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	JarID     *int64    `json:"jar_id"`
	JarName   string    `json:"jar_name"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"` // expense | income
	CreatedAt time.Time `json:"created_at"`
}
