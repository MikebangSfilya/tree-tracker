package plant

import (
	"context"
	"errors"
	"fmt"
)

type Store interface {
	GetCurrent(ctx context.Context) (Plant, error)
	Insert(ctx context.Context, p Plant) (Plant, error)
	Update(ctx context.Context, p Plant) (Plant, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetPlant(ctx context.Context) (Plant, error) {
	p, err := s.store.GetCurrent(ctx)
	if err != nil {
		return Plant{}, fmt.Errorf("get plant: %w", err)
	}
	return p, nil
}

// SavePlant persists the issued plant contract fields. If no row exists yet, it inserts one.
func (s *Service) SavePlant(ctx context.Context, p Plant) (Plant, error) {
	current, err := s.store.GetCurrent(ctx)
	if err != nil {
		if errors.Is(err, ErrPlantNotFound) {
			saved, insErr := s.store.Insert(ctx, p)
			if insErr != nil {
				return Plant{}, fmt.Errorf("save plant: %w", insErr)
			}
			return saved, nil
		}
		return Plant{}, fmt.Errorf("save plant: %w", err)
	}

	p.ID = current.ID
	saved, err := s.store.Update(ctx, p)
	if err != nil {
		return Plant{}, fmt.Errorf("save plant: %w", err)
	}
	return saved, nil
}

// ResetPlant restores the starting tree (epoch 0, phase 1, zero visual params).
// An existing seed is kept so the same plant identity is reused.
func (s *Service) ResetPlant(ctx context.Context) (Plant, error) {
	current, err := s.store.GetCurrent(ctx)
	initial := InitialPlant()
	if err != nil {
		if errors.Is(err, ErrPlantNotFound) {
			saved, insErr := s.store.Insert(ctx, initial)
			if insErr != nil {
				return Plant{}, fmt.Errorf("reset plant: %w", insErr)
			}
			return saved, nil
		}
		return Plant{}, fmt.Errorf("reset plant: %w", err)
	}

	initial.ID = current.ID
	initial.Seed = current.Seed
	saved, err := s.store.Update(ctx, initial)
	if err != nil {
		return Plant{}, fmt.Errorf("reset plant: %w", err)
	}
	return saved, nil
}
