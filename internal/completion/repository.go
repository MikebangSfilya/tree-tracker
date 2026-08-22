package completion

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type CompletionRepository1 struct {
	*pgx.Conn
	logger *slog.Logger
}

func NewCompletionRepository1(conn *pgx.Conn, logger slog.Logger) *CompletionRepository1 {
	return &CompletionRepository1{
		Conn:   conn,
		logger: &logger,
	}
}

func (s *CompletionRepository1) Complete(ctx context.Context, completion Completion) (Completion, error) {
	query := `
		INSERT INTO completions (user_id, routine_id, occurrence_key)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, routine_id, occurrence_key, completed_at;
	`

	var result Completion
	err := s.QueryRow(
		ctx,
		query,
		completion.UserID,
		completion.RoutineID,
		completion.OccurrenceKey).Scan(&result.ID,
		&result.UserID,
		&result.RoutineID,
		&result.OccurrenceKey,
		&result.CompletedAt)
	if err != nil {
		s.logger.Error("failed to insert completion", "err", err)
		return Completion{}, fmt.Errorf("insert completion: %w", err)
	}

	return result, nil
}
