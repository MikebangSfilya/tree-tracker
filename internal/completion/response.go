package completion

import (
	"time"

	"github.com/google/uuid"
)

type CompletionResponse struct {
	RoutineID   uuid.UUID `json:"routineID"`
	CompletedAt time.Time `json:"completedAt"`
}
