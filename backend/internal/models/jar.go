// Package models contains details about different tables
package models

import "time"

type Jar struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	AllocationType  string `json:"allocation_type"`
	AllocationValue int64  `json:"allocation_value"`
	SortOrder       int    `json:"sort_order"`
	IsSystem        bool   `json:"is_system"`

	CreatedAt time.Time `json:"created_at"`
}
