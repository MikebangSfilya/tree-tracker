package completion

import (
	"time"

	"github.com/google/uuid"
)

type Completion struct {
	Id            uuid.UUID `json:"id"`
	UserId        uuid.UUID `json:"userId"`
	RoutineId     int64     `json:"routineId"`
	OccurrenceKey string    `json:"occurrenceKey"`
	CompletedAt   time.Time `json:"completedAt"`
}

type CompletionRule struct {
	RoutineId  int64  `json:"routineId"`
	Recurrence string `json:"recurrence"`
	Active     bool   `json:"active"`
}
