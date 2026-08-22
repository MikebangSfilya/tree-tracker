package completion

import "github.com/google/uuid"

type CompleteInput struct {
	UserID    uuid.UUID
	RoutineID uuid.UUID
}

func NewCompleteInput(userID uuid.UUID, routineID uuid.UUID) CompleteInput {
	return CompleteInput{
		UserID:    userID,
		RoutineID: routineID,
	}
}
