package transactions

import "errors"

var (
	ErrTransactionNameRequired  = errors.New("name is required")
	ErrTransactionTypeRequired  = errors.New("type is required")
	ErrTransactionAmountInvalid = errors.New("amount must be greater than 0")
	ErrTransactionSameAccount   = errors.New("from and to account cannot be the same")
	ErrInvalidTransactionID     = errors.New("invalid transaction id")
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrJarNotFound              = errors.New("jar not found")
	ErrAccountNotFound          = errors.New("account not found")
)
