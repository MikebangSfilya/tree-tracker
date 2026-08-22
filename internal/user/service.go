package user

// TODO: 
// - login
// - access & refresh token
// - update
// - logout(?)

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


type RepoMethods interface {
	CreateUser(ctx context.Context, username, email string, password []byte) (id *uuid.UUID, err error)
	GetUser(ctx context.Context, email string) (user *User, err error)
	DeleteUser(ctx context.Context, email string) (err error)
}


type Service struct {
	logger *slog.Logger

	r RepoMethods
}


func NewService(logger *slog.Logger, r RepoMethods) *Service {
	return &Service{
		logger: logger,
		r: r,
	}
}


func (s Service) Register(ctx context.Context, username, email, password string) (string, error) {
	// encrypt parol
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(
		"failed to hash password",
		"password string", password,
		"error", err,
		)
		return "", err
	}

	id, err := s.r.CreateUser(ctx, username, email, hashedPassword)
	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return "", err
	}

	return id.String(), nil
}

// not implemented
func (s Service) Login(ctx context.Context, email, password string) (string, error) {
	return "", nil
}


func (s Service) GetUser(ctx context.Context, email string) (*User, error) {
	user, err := s.r.GetUser(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user", "error", err)
		return nil, err
	}

	return user, nil
}


func (s Service) DeleteUser(ctx context.Context, email string) error {
	if err := s.r.DeleteUser(ctx, email); err != nil {
		s.logger.Error("failed to delete user", "error", err)
		return err
	}

	return nil
}
