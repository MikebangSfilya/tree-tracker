package user

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type stubRepository struct {
	createUser func(context.Context, string, string, []byte) (*uuid.UUID, error)
	getUser    func(context.Context, string) (*User, error)
	deleteUser func(context.Context, string) error
}

func (r stubRepository) CreateUser(ctx context.Context, username, email string, password []byte) (*uuid.UUID, error) {
	return r.createUser(ctx, username, email, password)
}

func (r stubRepository) GetUser(ctx context.Context, email string) (*User, error) {
	return r.getUser(ctx, email)
}

func (r stubRepository) DeleteUser(ctx context.Context, email string) error {
	return r.deleteUser(ctx, email)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestServiceRegister(t *testing.T) {
	t.Parallel()

	wantID := uuid.New()
	repo := stubRepository{
		createUser: func(_ context.Context, username, email string, password []byte) (*uuid.UUID, error) {
			if username != "mike" {
				t.Fatalf("username = %q, want mike", username)
			}
			if email != "mike@example.com" {
				t.Fatalf("email = %q, want mike@example.com", email)
			}
			if err := bcrypt.CompareHashAndPassword(password, []byte("secret")); err != nil {
				t.Fatalf("stored password is not a valid hash: %v", err)
			}
			return &wantID, nil
		},
	}

	service := NewService(testLogger(), repo)
	gotID, err := service.Register(context.Background(), "mike", "mike@example.com", "secret")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if gotID != wantID.String() {
		t.Fatalf("Register() id = %q, want %q", gotID, wantID)
	}
}

func TestServiceRegisterReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	repo := stubRepository{
		createUser: func(context.Context, string, string, []byte) (*uuid.UUID, error) {
			return nil, ErrUserAlreadyExists
		},
	}

	service := NewService(testLogger(), repo)
	_, err := service.Register(context.Background(), "mike", "mike@example.com", "secret")
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf("Register() error = %v, want %v", err, ErrUserAlreadyExists)
	}
}

func TestServiceGetUser(t *testing.T) {
	t.Parallel()

	want := &User{ID: uuid.New(), Username: "mike", Email: "mike@example.com"}
	repo := stubRepository{
		getUser: func(_ context.Context, email string) (*User, error) {
			if email != want.Email {
				t.Fatalf("email = %q, want %q", email, want.Email)
			}
			return want, nil
		},
	}

	service := NewService(testLogger(), repo)
	got, err := service.GetUser(context.Background(), want.Email)
	if err != nil {
		t.Fatalf("GetUser() error = %v", err)
	}
	if got != want {
		t.Fatalf("GetUser() = %#v, want %#v", got, want)
	}
}

func TestServiceDeleteUserReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	wantErr := errors.New("delete failed")
	repo := stubRepository{
		deleteUser: func(_ context.Context, email string) error {
			if email != "mike@example.com" {
				t.Fatalf("email = %q, want mike@example.com", email)
			}
			return wantErr
		},
	}

	service := NewService(testLogger(), repo)
	err := service.DeleteUser(context.Background(), "mike@example.com")
	if !errors.Is(err, wantErr) {
		t.Fatalf("DeleteUser() error = %v, want %v", err, wantErr)
	}
}
