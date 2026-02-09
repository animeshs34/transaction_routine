package errors

import "errors"

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidAmount   = errors.New("invalid amount")
	ErrOperationType   = errors.New("invalid operation type")
	ErrInternal        = errors.New("internal server error")
)
