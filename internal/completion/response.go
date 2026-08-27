package completion

import (
	"time"
)

type CompletionResponse struct {
	Id          int64     `json:"id"`
	RoutineID   int64     `json:"routineId"`
	CompletedAt time.Time `json:"completedAt"`
}
