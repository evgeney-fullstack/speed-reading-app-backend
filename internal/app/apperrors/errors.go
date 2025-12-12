package apperrors

import "errors"

var (
	ErrTextNotFound = errors.New("text not found")
	ErrInvalidID    = errors.New("invalid ID")
)
