package jars

import "errors"

var (
	ErrJarNameRequired          = errors.New("jar name required")
	ErrInvalidJarID             = errors.New("invalid jar id")
	ErrPercentageMustBePositive = errors.New("percentage must be positive")
	ErrTotalPercentageExceeded  = errors.New("total percentage exceeds 100")
	ErrJarNameAlreadyExists     = errors.New("jar name already exists")
	ErrInvalidAllocationType    = errors.New("invalid or missing allocation type")
	ErrJarNotFound              = errors.New("jar not found")
	ErrJarValidation            = errors.New("jar validation failed")
)
