// Package jars
package jars

import "time"

type AllocationType string

const (
	AllocationPercentage AllocationType = "percentage"
	AllocationRemainder  AllocationType = "remainder"
)

type Jar struct {
	ID             int64          `json:"id"`
	Name           string         `json:"name"`
	AllocationType AllocationType `json:"allocation_type"`

	// Used only when AllocationType == percentage
	Value int64 `json:"value"`

	CreatedAt time.Time `json:"created_at"`
}

type JarWithAllocation struct {
	Jar
	AllocatedAmount int64 `json:"allocated_amount"`
}
