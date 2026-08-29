package completion

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type CompletionRepository1 struct {
	q      Querier
	logger *slog.Logger
}

func NewCompletionRepository1(q Querier, logger *slog.Logger) *CompletionRepository1 {
	return &CompletionRepository1{
		q:      q,
		logger: logger,
	}
}

func (s *CompletionRepository1) Complete(ctx context.Context, completion Completion) (Completion, error) {
	query := `
		INSERT INTO completions (user_id, routine_id, occurrence_key)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, routine_id, occurrence_key, completed_at;
	`

	var result Completion
	err := s.q.QueryRow(
		ctx,
		query,
		completion.UserId,
		completion.RoutineId,
		completion.OccurrenceKey).Scan(&result.Id,
		&result.UserId,
		&result.RoutineId,
		&result.OccurrenceKey,
		&result.CompletedAt)
	if err != nil {
		s.logger.Error("failed to insert completion", "err", err)
		return Completion{}, fmt.Errorf("insert completion: %w", err)
	}

	return result, nil
}

func (s *CompletionRepository1) GetCompletions(ctx context.Context, userID uuid.UUID) ([]Completion, error) {
	query := `
		SELECT id, user_id, routine_id, occurrence_key, completed_at FROM completions WHERE user_id = $1;
`
	rows, err := s.q.Query(ctx, query, userID)
	if err != nil {
		s.logger.Error("failed to get completions", "err", err)
		return nil, fmt.Errorf("get completions: %w", err)
	}
	defer rows.Close()
	var completions []Completion
	for rows.Next() {
		var completion Completion
		if err = rows.Scan(&completion.Id,
			&completion.UserId,
			&completion.RoutineId,
			&completion.OccurrenceKey,
			&completion.CompletedAt); err != nil {
			s.logger.Error("failed to scan completions", "err", err)
			return nil, fmt.Errorf("get completions: %w", err)
		}
		completions = append(completions, completion)
	}

	return completions, nil

}
