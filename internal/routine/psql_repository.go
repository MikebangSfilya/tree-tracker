package routine

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type postgresRepository struct {
	q Querier
}

func NewPostgresRepository(q Querier) Repository {
	return &postgresRepository{
		q: q,
	}
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]Routine, error) {
	const query = `
		SELECT
			id,
			category,
			type,
			weight,
			coefficient,
			temporary,
			disposable,
			created_at,
			updated_at
		FROM routines
		ORDER BY id
	`

	rows, err := r.q.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	routines := make([]Routine, 0)

	for rows.Next() {
		var routine Routine

		err := rows.Scan(
			&routine.ID,
			&routine.Category,
			&routine.Type,
			&routine.Weight,
			&routine.Coefficient,
			&routine.Temporary,
			&routine.Disposable,
			&routine.CreatedAt,
			&routine.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		routines = append(routines, routine)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routines, nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int64) (Routine, error) {
	const query = `
		SELECT
			id,
			category,
			type,
			weight,
			coefficient,
			temporary,
			disposable,
			created_at,
			updated_at
		FROM routines
		WHERE id = $1
	`

	var routine Routine

	err := r.q.QueryRow(ctx, query, id).Scan(
		&routine.ID,
		&routine.Category,
		&routine.Type,
		&routine.Weight,
		&routine.Coefficient,
		&routine.Temporary,
		&routine.Disposable,
		&routine.CreatedAt,
		&routine.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Routine{}, ErrRoutineNotFound
		}

		return Routine{}, err
	}

	return routine, nil
}

func (r *postgresRepository) Create(
	ctx context.Context,
	routine Routine,
) (Routine, error) {
	const query = `
		INSERT INTO routines (
			category,
			type,
			weight,
			coefficient,
			temporary,
			disposable
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			category,
			type,
			weight,
			coefficient,
			temporary,
			disposable,
			created_at,
			updated_at
	`

	var created Routine

	err := r.q.QueryRow(
		ctx,
		query,
		routine.Category,
		routine.Type,
		routine.Weight,
		routine.Coefficient,
		routine.Temporary,
		routine.Disposable,
	).Scan(
		&created.ID,
		&created.Category,
		&created.Type,
		&created.Weight,
		&created.Coefficient,
		&created.Temporary,
		&created.Disposable,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return Routine{}, err
	}

	return created, nil
}

func (r *postgresRepository) Update(
	ctx context.Context,
	id int64,
	request UpdateRequest,
) (Routine, error) {
	const query = `
		UPDATE routines
		SET
			category = COALESCE($2, category),
			type = COALESCE($3, type),
			weight = COALESCE($4, weight),
			coefficient = COALESCE($5, coefficient),
			temporary = COALESCE($6, temporary),
			disposable = COALESCE($7, disposable),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING
			id,
			category,
			type,
			weight,
			coefficient,
			temporary,
			disposable,
			created_at,
			updated_at
	`

	var routine Routine

	err := r.q.QueryRow(
		ctx,
		query,
		id,
		request.Category,
		request.Type,
		request.Weight,
		request.Coefficient,
		request.Temporary,
		request.Disposable,
	).Scan(
		&routine.ID,
		&routine.Category,
		&routine.Type,
		&routine.Weight,
		&routine.Coefficient,
		&routine.Temporary,
		&routine.Disposable,
		&routine.CreatedAt,
		&routine.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Routine{}, ErrRoutineNotFound
		}

		return Routine{}, err
	}

	return routine, nil
}

func (r *postgresRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM routines
		WHERE id = $1
	`

	result, err := r.q.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrRoutineNotFound
	}

	return nil
}
