package models

import "errors"

var (
	ErrNotFound         = errors.New("resource not found")
	ErrUnauthorized     = errors.New("authentication required")
	ErrForbidden        = errors.New("permission denied")
	ErrInvalidInput     = errors.New("invalid input")
	ErrConflict         = errors.New("resource already exists")
	ErrInternal         = errors.New("internal server error")
	ErrPredictionLocked = errors.New("prediction period has closed")
	ErrSessionNotFound  = errors.New("requested session not found")
)
