package progress

import (
	"context"

	"github.com/MikebangSfilya/tree-tracker/internal/completion"
	"github.com/google/uuid"
)

type Repository interface {
	CompleteAndGrow(ctx context.Context, userID uuid.UUID, routineID int64) (completion.Completion, error)
}

type CompletionReader interface {
	GetCompletions(ctx context.Context, userID uuid.UUID) ([]completion.Completion, error)
}

type Service struct {
	repository Repository
	reader     CompletionReader
}

func NewService(repository Repository, reader CompletionReader) *Service {
	return &Service{repository: repository, reader: reader}
}

func (s *Service) RoutineComplete(ctx context.Context, input completion.CompleteInput) (completion.Completion, error) {
	return s.repository.CompleteAndGrow(ctx, input.UserId, input.RoutineId)
}

func (s *Service) GetCompletions(ctx context.Context, userID uuid.UUID) ([]completion.Completion, error) {
	return s.reader.GetCompletions(ctx, userID)
}
