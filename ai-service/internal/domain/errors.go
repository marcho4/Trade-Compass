package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrUnknownTaskType    = errors.New("unknown task type")
	ErrScenariosNotFound  = errors.New("scenarios not found")
	ErrDCFResultsNotFound = errors.New("dcf results not found")
)
