package models

import "time"

type Jar struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	AllocationValue float64   `json:"allocation_value"`
	SortOrder       int       `json:"sort_order"`
	IsSystem        bool      `json:"is_system"`
	CreatedAt       time.Time `json:"created_at"`
}
