package accounts

import "errors"

var (
	ErrAccountAlreadyExists   = errors.New("account already exists")
	ErrInvalidAccountID       = errors.New("invalid account id")
	ErrInvalidAccountName     = errors.New("invalid account name")
	ErrAccountAlreadyArchived = errors.New("account already archived")
	ErrAccountAlreadyActive   = errors.New("account already active")
	ErrAccountNotFound        = errors.New("account not found")
)
