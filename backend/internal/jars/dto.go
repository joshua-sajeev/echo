package jars

type CreateJarRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=100"`
	AllocationType string  `json:"allocation_type" validate:"required,oneof=percentage fixed_amount remainder"`
	Value          float64 `json:"value" validate:"gte=0,lte=100"`
	Priority       int     `json:"priority" validate:"gte=0"`
}

type UpdateJarRequest struct {
	Name           string  `json:"name" validate:"required,min=1,max=100"`
	AllocationType string  `json:"allocation_type" validate:"required,oneof=percentage fixed_amount remainder"`
	Value          float64 `json:"value" validate:"gte=0,lte=100"`
	Priority       int     `json:"priority" validate:"gte=0"`
}
