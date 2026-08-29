package progress

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MikebangSfilya/tree-tracker/internal/completion"
	"github.com/MikebangSfilya/tree-tracker/internal/plant"
	completionErrors "github.com/MikebangSfilya/tree-tracker/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryDB struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *RepositoryDB {
	return &RepositoryDB{db: db}
}

func (r *RepositoryDB) CompleteAndGrow(ctx context.Context, userID uuid.UUID, routineID int64) (completion.Completion, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return completion.Completion{}, fmt.Errorf("begin progress transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var routine struct {
		weight      int
		coefficient int
		disposable  bool
	}
	err = tx.QueryRow(ctx, "SELECT weight, coefficient, disposable FROM routines WHERE id = $1", routineID).Scan(&routine.weight, &routine.coefficient, &routine.disposable)
	if errors.Is(err, pgx.ErrNoRows) {
		return completion.Completion{}, fmt.Errorf("routine %d not found", routineID)
	}
	if err != nil {
		return completion.Completion{}, fmt.Errorf("get routine for progress: %w", err)
	}

	occurrenceKey := time.Now().UTC().Format("2006-01-02")
	if routine.disposable {
		occurrenceKey = "once"
	}
	var result completion.Completion
	err = tx.QueryRow(ctx, `
		INSERT INTO completions (user_id, routine_id, occurrence_key)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, routine_id, occurrence_key, completed_at`, userID, routineID, occurrenceKey,
	).Scan(&result.Id, &result.UserId, &result.RoutineId, &result.OccurrenceKey, &result.CompletedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if routine.disposable {
				return completion.Completion{}, completionErrors.ErrAlreadyCompleted
			}
			return completion.Completion{}, completionErrors.ErrAlreadyCompletedToday
		}
		return completion.Completion{}, fmt.Errorf("save completion: %w", err)
	}

	current, err := currentPlant(ctx, tx)
	if err != nil {
		return completion.Completion{}, err
	}
	next := Advance(current, routine.weight*routine.coefficient)
	if _, err := tx.Exec(ctx, "UPDATE plants SET epoch = $2, phase = $3, phase_progress = $4 WHERE id = $1", next.ID, next.Epoch, next.Phase, next.PhaseProgress); err != nil {
		return completion.Completion{}, fmt.Errorf("grow plant: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return completion.Completion{}, fmt.Errorf("commit progress transaction: %w", err)
	}
	return result, nil
}

func currentPlant(ctx context.Context, tx pgx.Tx) (plant.Plant, error) {
	var p plant.Plant
	err := tx.QueryRow(ctx, `
		SELECT id, epoch, phase, phase_progress, branching, density, curvature, vitality, seed
		FROM plants ORDER BY id DESC LIMIT 1 FOR UPDATE`,
	).Scan(&p.ID, &p.Epoch, &p.Phase, &p.PhaseProgress, &p.Branching, &p.Density, &p.Curvature, &p.Vitality, &p.Seed)
	if errors.Is(err, pgx.ErrNoRows) {
		return plant.Plant{}, errors.New("plant not found")
	}
	if err != nil {
		return plant.Plant{}, fmt.Errorf("get plant for progress: %w", err)
	}
	return p, nil
}
