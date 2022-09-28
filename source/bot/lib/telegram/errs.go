package telegram

import "errors"

var (
	ErrUnknown      = errors.New("unknown error")
	ErrUnauthorized = errors.New("not authorized")
)
