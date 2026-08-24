package completion

import "github.com/google/uuid"

type CompleteInput struct {
	UserId    uuid.UUID
	RoutineId int64
}

func NewCompleteInput(userId uuid.UUID, routineId int64) CompleteInput {
	return CompleteInput{
		UserId:    userId,
		RoutineId: routineId,
	}
}
