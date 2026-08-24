package plant

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Querier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Repository struct {
	db Querier
}

func NewRepository(db Querier) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetCurrent(ctx context.Context) (Plant, error) {
	const q = `
		SELECT id, epoch, phase, phase_progress, branching, density, curvature, vitality, seed
		FROM plants
		ORDER BY id DESC
		LIMIT 1`

	var p Plant
	err := r.db.QueryRow(ctx, q).Scan(
		&p.ID,
		&p.Epoch,
		&p.Phase,
		&p.PhaseProgress,
		&p.Branching,
		&p.Density,
		&p.Curvature,
		&p.Vitality,
		&p.Seed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Plant{}, ErrPlantNotFound
		}
		return Plant{}, fmt.Errorf("get current plant: %w", err)
	}

	return p, nil
}
