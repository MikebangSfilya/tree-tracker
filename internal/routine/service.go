package routine

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidName        = errors.New("invalid routine name")
	ErrInvalidDescription = errors.New("invalid routine description")
	ErrInvalidCategory    = errors.New("invalid routine category")
	ErrInvalidWeight      = errors.New("invalid routine weight")
	ErrInvalidCoefficient = errors.New("invalid routine coefficient")
	ErrRoutineNotFound    = errors.New("routine not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]Routine, error)
	Create(ctx context.Context, request CreateRequest) (Routine, error)
	Update(ctx context.Context, id int64, request UpdateRequest) (Routine, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) GetAll(ctx context.Context) ([]Routine, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) Create(ctx context.Context, request CreateRequest) (Routine, error) {
	request.Name = strings.TrimSpace(request.Name)
	request.Description = strings.TrimSpace(request.Description)
	request.Category = strings.ToLower(strings.TrimSpace(request.Category))

	if !isValidName(request.Name) {
		return Routine{}, ErrInvalidName
	}

	if !isValidDescription(request.Description) {
		return Routine{}, ErrInvalidDescription
	}

	if !isValidCategory(request.Category) {
		return Routine{}, ErrInvalidCategory
	}

	if !isValidWeight(request.Weight) {
		return Routine{}, ErrInvalidWeight
	}

	if !isValidCoefficient(request.Coefficient) {
		return Routine{}, ErrInvalidCoefficient
	}

	routine := Routine{
		Name:        request.Name,
		Description: request.Description,
		Category:    request.Category,
		Weight:      request.Weight,
		Coefficient: request.Coefficient,
	}

	return s.repository.Create(ctx, routine)
}

func (s *service) Update(ctx context.Context, id int64, request UpdateRequest) (Routine, error) {
	if id <= 0 {
		return Routine{}, ErrRoutineNotFound
	}

	routine, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Routine{}, err
	}

	if request.Name != nil {
		routine.Name = strings.TrimSpace(*request.Name)
	}

	if request.Description != nil {
		routine.Description = strings.TrimSpace(*request.Description)
	}

	if request.Category != nil {
		routine.Category = strings.ToLower(strings.TrimSpace(*request.Category))
	}

	if request.Weight != nil {
		routine.Weight = *request.Weight
	}

	if request.Coefficient != nil {
		routine.Coefficient = *request.Coefficient
	}

	if !isValidName(routine.Name) {
		return Routine{}, ErrInvalidName
	}

	if !isValidDescription(routine.Description) {
		return Routine{}, ErrInvalidDescription
	}

	if !isValidCategory(routine.Category) {
		return Routine{}, ErrInvalidCategory
	}

	if !isValidWeight(routine.Weight) {
		return Routine{}, ErrInvalidWeight
	}

	if !isValidCoefficient(routine.Coefficient) {
		return Routine{}, ErrInvalidCoefficient
	}

	return s.repository.Update(ctx, routine)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrRoutineNotFound
	}

	return s.repository.Delete(ctx, id)
}

func isValidName(name string) bool {
	return name != "" && len(name) <= 100
}

func isValidDescription(description string) bool {
	return len(description) <= 500
}

func isValidCategory(category string) bool {
	switch category {
	case "health", "study", "personal":
		return true
	default:
		return false
	}
}

func isValidWeight(weight int) bool {
	return weight >= 1 && weight <= 5
}

func isValidCoefficient(coefficient int) bool {
	return coefficient >= 1 && coefficient <= 5
}
