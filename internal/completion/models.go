package completion

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAlreadyCompleted      = errors.New("routine is already completed")
	ErrAlreadyCompletedToday = errors.New("routine had already been completed today")
)

type Completion struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"userId"`
	RoutineID     uuid.UUID `json:"routineId"`
	OccurrenceKey string    `json:"occurrenceKey"`
	CompletedAt   time.Time `json:"completedAt"`
}

type CompletionRule struct {
	RoutineID  uuid.UUID
	Recurrence string
	Active     bool
}
