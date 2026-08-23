package completion

import (
	"time"

	"github.com/google/uuid"
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
