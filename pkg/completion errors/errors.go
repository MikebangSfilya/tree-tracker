package completion_errors

import "errors"

var (
	ErrAlreadyCompleted      = errors.New("routine is already completed")
	ErrAlreadyCompletedToday = errors.New("routine had already been completed today")
)
