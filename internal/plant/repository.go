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

const plantColumns = `id, epoch, phase, phase_progress, branching, density, curvature, vitality, seed`

func (r *Repository) GetCurrent(ctx context.Context) (Plant, error) {
	const q = `
		SELECT ` + plantColumns + `
		FROM plants
		ORDER BY id DESC
		LIMIT 1`

	p, err := scanPlant(r.db.QueryRow(ctx, q))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Plant{}, ErrPlantNotFound
		}
		return Plant{}, fmt.Errorf("get current plant: %w", err)
	}

	return p, nil
}

func (r *Repository) Insert(ctx context.Context, p Plant) (Plant, error) {
	const q = `
		INSERT INTO plants (epoch, phase, phase_progress, branching, density, curvature, vitality, seed)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING ` + plantColumns

	saved, err := scanPlant(r.db.QueryRow(ctx, q,
		p.Epoch,
		p.Phase,
		p.PhaseProgress,
		p.Branching,
		p.Density,
		p.Curvature,
		p.Vitality,
		p.Seed,
	))
	if err != nil {
		return Plant{}, fmt.Errorf("insert plant: %w", err)
	}

	return saved, nil
}

func (r *Repository) Update(ctx context.Context, p Plant) (Plant, error) {
	const q = `
		UPDATE plants
		SET epoch = $1,
		    phase = $2,
		    phase_progress = $3,
		    branching = $4,
		    density = $5,
		    curvature = $6,
		    vitality = $7,
		    seed = $8
		WHERE id = $9
		RETURNING ` + plantColumns

	saved, err := scanPlant(r.db.QueryRow(ctx, q,
		p.Epoch,
		p.Phase,
		p.PhaseProgress,
		p.Branching,
		p.Density,
		p.Curvature,
		p.Vitality,
		p.Seed,
		p.ID,
	))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Plant{}, ErrPlantNotFound
		}
		return Plant{}, fmt.Errorf("update plant: %w", err)
	}

	return saved, nil
}

func scanPlant(row pgx.Row) (Plant, error) {
	var p Plant
	err := row.Scan(
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
		return Plant{}, err
	}
	return p, nil
}
