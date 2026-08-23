package completion

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	completion_errors "github.com/MikebangSfilya/tree-tracker/pkg/completion errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type CompletionService1 struct {
	repository     CompletionRepository
	routineFetcher RoutineFetcher
	logger         *slog.Logger
}

type CompletionRepository interface {
	Complete(ctx context.Context, completion Completion) (Completion, error)
	GetCompletions(ctx context.Context, userID uuid.UUID) ([]Completion, error)
}

type RoutineFetcher interface {
	GetCompletionRule(ctx context.Context, uuid uuid.UUID) (*CompletionRule, error)
}

func NewCompletionService1(repository CompletionRepository, routineFetcher RoutineFetcher, logger slog.Logger) *CompletionService1 {
	return &CompletionService1{
		repository:     repository,
		routineFetcher: routineFetcher,
		logger:         &logger}
}

func (s *CompletionService1) RoutineComplete(ctx context.Context, input CompleteInput) (Completion, error) {
	rule, err := s.routineFetcher.GetCompletionRule(ctx, input.RoutineID)
	if err != nil {
		s.logger.Error("failed to get completion rule")
		return Completion{}, fmt.Errorf("failed to get completion rule: %w", err)
	}

	if !rule.Active {
		s.logger.Warn("completion rule is not active")
		return Completion{}, fmt.Errorf("active already completed: %w", completion_errors.ErrAlreadyCompleted)
	}

	occurrenceKey := generateOccurrenceKey(rule.Recurrence, time.Now().UTC())
	s.logger.Debug("completion rule occurred", slog.String("key", occurrenceKey))

	completion := Completion{
		UserID:        input.UserID,
		RoutineID:     input.RoutineID,
		OccurrenceKey: occurrenceKey,
	}
	result, err := s.repository.Complete(ctx, completion)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if rule.Recurrence == "once" {
				return Completion{}, fmt.Errorf("completion rule already completed: %w", completion_errors.ErrAlreadyCompleted)
			}
			if rule.Recurrence == "daily" {
				return Completion{}, fmt.Errorf("completion rule already completed: %w", completion_errors.ErrAlreadyCompletedToday)
			}
		}

		s.logger.Error("failed to save completion", "err", err)
		return Completion{}, fmt.Errorf("failed to save completion: %w", err)
	}
	return result, nil
}

func (s *CompletionService1) GetCompletions(ctx context.Context, userID uuid.UUID) ([]Completion, error) {
	completions, err := s.repository.GetCompletions(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get completions")
		return nil, fmt.Errorf("failed to get completions: %w", err)
	}
	return completions, nil
}

func generateOccurrenceKey(recurrence string, now time.Time) string {
	switch recurrence {
	case "daily":
		return now.Format("2006-01-02")
	case "once":
		return "once"
	default:
		return "default"
	}
}
