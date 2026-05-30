package accounts

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyState  = errors.New("account already in requested state")
	ErrInvalidAccountID     = errors.New("invalid account id")
	ErrInvalidAccountName   = errors.New("invalid account name")
)
