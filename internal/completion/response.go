package completion

import (
	"time"
)

type CompletionResponse struct {
	RoutineID   int64     `json:"routineId"`
	CompletedAt time.Time `json:"completedAt"`
}
