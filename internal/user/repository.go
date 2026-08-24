package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)


type RepoDB struct {
	db *sql.DB
}


func NewRepoDB(db *sql.DB) *RepoDB {
	return &RepoDB{db: db}
}


func (r RepoDB) CreateUser(ctx context.Context, username, email string, password []byte) (*uuid.UUID, error) {
	q := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id uuid.UUID
	err := r.db.QueryRowContext(ctx, q, username, email, password).Scan(&id)
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
	err := r.db.QueryRowContext(ctx, q, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		// check esli takogo usera net
		if err == pgx.ErrNoRows {
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

	res, err := r.db.ExecContext(ctx, q, email);
	if err != nil {
		return err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAff < 1 {
		return ErrFailedToDeleteUser
	}

	return nil
}
