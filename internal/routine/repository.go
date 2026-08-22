package routine

import "context"

type Repository interface {
	GetAll(ctx context.Context) ([]Routine, error)
	GetByID(ctx context.Context, id int64) (Routine, error)
	Create(ctx context.Context, routine Routine) (Routine, error)
	Update(ctx context.Context, routine Routine) (Routine, error)
	Delete(ctx context.Context, id int64) error
}
