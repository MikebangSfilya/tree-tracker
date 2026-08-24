package plant

import (
	"context"
	"fmt"
)

type Store interface {
	GetCurrent(ctx context.Context) (Plant, error)
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
