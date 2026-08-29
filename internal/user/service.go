package user

// TODO:
// - login
// - access & refresh token
// - update
// - logout(?)

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	demoUsername = "Demo user"
	demoEmail    = "demo@tree-tracker.local"
	demoPassword = "demo-tree-2026"
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
		r:      r,
	}
}

func (s Service) Register(ctx context.Context, username, email, password string) (string, error) {
	// encrypt parol
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(
			"failed to hash password",
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

func (s Service) EnsureDemoUser(ctx context.Context) (string, error) {
	existing, err := s.r.GetUser(ctx, demoEmail)
	if err == nil {
		return existing.ID.String(), nil
	}
	if !errors.Is(err, ErrUserNotFound) {
		s.logger.Error("failed to find demo user", "error", err)
		return "", err
	}

	id, err := s.Register(ctx, demoUsername, demoEmail, demoPassword)
	if !errors.Is(err, ErrUserAlreadyExists) {
		return id, err
	}

	existing, err = s.r.GetUser(ctx, demoEmail)
	if err != nil {
		return "", err
	}
	return existing.ID.String(), nil
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
