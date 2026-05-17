// Package models
package models

import "time"

type Jar struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	AllocationType  string  `json:"allocation_type"`
	AllocationValue float64 `json:"allocation_value"`

	RoundTo   int  `json:"round_to"`
	SortOrder int  `json:"sort_order"`
	IsActive  bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
