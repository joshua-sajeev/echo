// Package jars
package jars

import "time"

type AllocationType string

const (
	AllocationPercentage AllocationType = "percentage"
	AllocationFixed      AllocationType = "fixed_amount"
	AllocationRemainder  AllocationType = "remainder"
)

type Jar struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	AllocationType AllocationType `json:"allocation_type"`
	Value          int64          `json:"value"`
	Priority       int            `json:"priority"`
	CreatedAt      time.Time      `json:"created_at"`
}
