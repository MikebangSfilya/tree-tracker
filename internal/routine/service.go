package routine

import (
	"context"
	"errors"
	"strings"

	"github.com/MikebangSfilya/tree-tracker/internal/completion"
)

var (
	ErrInvalidCategory    = errors.New("invalid routine category")
	ErrInvalidType        = errors.New("invalid routine type")
	ErrInvalidWeight      = errors.New("invalid routine weight")
	ErrInvalidCoefficient = errors.New("invalid routine coefficient")
	ErrRoutineNotFound    = errors.New("routine not found")
)

type Service interface {
	GetAll(ctx context.Context) ([]Routine, error)
	Create(ctx context.Context, request CreateRequest) (Routine, error)
	Update(ctx context.Context, id int64, request UpdateRequest) (Routine, error)
	Delete(ctx context.Context, id int64) error
	GetCompletionRule(ctx context.Context, routineID int64) (*completion.CompletionRule, error)
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
	request.Category = strings.ToLower(strings.TrimSpace(request.Category))
	request.Type = strings.TrimSpace(request.Type)

	if !isValidCategory(request.Category) {
		return Routine{}, ErrInvalidCategory
	}

	if !isValidType(request.Type) {
		return Routine{}, ErrInvalidType
	}

	if !isValidWeight(request.Weight) {
		return Routine{}, ErrInvalidWeight
	}

	if !isValidCoefficient(request.Coefficient) {
		return Routine{}, ErrInvalidCoefficient
	}

	routine := Routine{
		Category:    request.Category,
		Type:        request.Type,
		Weight:      request.Weight,
		Coefficient: request.Coefficient,
		Temporary:   request.Temporary,
		Disposable:  request.Disposable,
	}

	return s.repository.Create(ctx, routine)
}

func (s *service) Update(
	ctx context.Context,
	id int64,
	request UpdateRequest,
) (Routine, error) {
	if id <= 0 {
		return Routine{}, ErrRoutineNotFound
	}

	routine, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Routine{}, err
	}

	if request.Category != nil {
		routine.Category = strings.ToLower(strings.TrimSpace(*request.Category))
	}

	if request.Type != nil {
		routine.Type = strings.TrimSpace(*request.Type)
	}

	if request.Weight != nil {
		routine.Weight = *request.Weight
	}

	if request.Coefficient != nil {
		routine.Coefficient = *request.Coefficient
	}

	if request.Temporary != nil {
		routine.Temporary = *request.Temporary
	}

	if request.Disposable != nil {
		routine.Disposable = *request.Disposable
	}

	if !isValidCategory(routine.Category) {
		return Routine{}, ErrInvalidCategory
	}

	if !isValidType(routine.Type) {
		return Routine{}, ErrInvalidType
	}

	if !isValidWeight(routine.Weight) {
		return Routine{}, ErrInvalidWeight
	}

	if !isValidCoefficient(routine.Coefficient) {
		return Routine{}, ErrInvalidCoefficient
	}

	return s.repository.Update(ctx, id, UpdateRequest{
		Category:    &routine.Category,
		Type:        &routine.Type,
		Weight:      &routine.Weight,
		Coefficient: &routine.Coefficient,
		Temporary:   &routine.Temporary,
		Disposable:  &routine.Disposable,
	})
}

func (s *service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrRoutineNotFound
	}

	return s.repository.Delete(ctx, id)
}

func isValidCategory(category string) bool {
	switch category {
	case "health", "study", "personal":
		return true
	default:
		return false
	}
}

func isValidType(routineType string) bool {
	return routineType != "" && len(routineType) <= 100
}

func isValidWeight(weight int) bool {
	return weight >= 1 && weight <= 5
}

func isValidCoefficient(coefficient int) bool {
	return coefficient >= 1 && coefficient <= 5
}

func (s *service) GetCompletionRule(
	ctx context.Context,
	routineID int64,
) (*completion.CompletionRule, error) {
	routine, err := s.repository.GetByID(ctx, routineID)
	if err != nil {
		return nil, err
	}

	recurrence := "daily"

	if routine.Disposable {
		recurrence = "once"
	}

	return &completion.CompletionRule{
		RoutineId:  routine.ID,
		Recurrence: recurrence,
		Active:     true,
	}, nil
}
