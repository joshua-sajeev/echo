package auth

import "errors"

var ErrInvalidPIN = errors.New(
	"invalid pin",
)
