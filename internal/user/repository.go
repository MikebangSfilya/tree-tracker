package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type RepoDB struct {
	q Querier
}

func NewRepoDB(q Querier) *RepoDB {
	return &RepoDB{q}
}

func (r RepoDB) CreateUser(ctx context.Context, username, email string, password []byte) (*uuid.UUID, error) {
	q := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id uuid.UUID
	err := r.q.QueryRow(ctx, q, username, email, password).Scan(&id)
	if err != nil {
		// check esli takoy email uje zaregan
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return nil, ErrUserAlreadyExists
		}

		return nil, err
	}

	return &id, nil
}

func (r RepoDB) GetUser(ctx context.Context, email string) (*User, error) {
	q := `
		SELECT id, username, email, password, created_at
		FROM users
		WHERE email = $1;
	`

	var user User
	err := r.q.QueryRow(ctx, q, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		// check esli takogo usera net
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r RepoDB) DeleteUser(ctx context.Context, email string) error {
	q := `
		DELETE FROM users
		WHERE email = $1;
	`

	res, err := r.q.Exec(ctx, q, email)
	if err != nil {
		return err
	}

	rowsAff := res.RowsAffected()
	if rowsAff < 1 {
		return ErrUserNotFound
	}

	return nil
}
