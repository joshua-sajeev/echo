package jars

type CreateJarRequest struct {
	Name           string `json:"name" validate:"required,min=1,max=100"`
	AllocationType string `json:"allocation_type" validate:"required,oneof=percentage remainder"`
	// Percentage value (0 for remainder jars)
	Value int64 `json:"value" validate:"gte=0,lte=100"`
}

type UpdateJarRequest struct {
	Name           *string `json:"name" validate:"omitempty,min=1,max=100"`
	AllocationType *string `json:"allocation_type" validate:"omitempty,oneof=percentage remainder"`
	Value          *int64  `json:"value" validate:"omitempty,gte=0,lte=100"`
}
