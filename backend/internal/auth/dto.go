package auth

type LoginRequest struct {
	PIN string `json:"pin" validate:"required"`
}
