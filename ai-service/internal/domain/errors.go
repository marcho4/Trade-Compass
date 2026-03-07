package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrUnknownTaskType = errors.New("unknown task type")
)
