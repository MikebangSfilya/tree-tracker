package plant

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
)

type memoryStore struct {
	plant  *Plant
	err    error
	nextID int64
}

func (m *memoryStore) GetCurrent(context.Context) (Plant, error) {
	if m.err != nil {
		return Plant{}, m.err
	}
	if m.plant == nil {
		return Plant{}, ErrPlantNotFound
	}
	return *m.plant, nil
}

func (m *memoryStore) Insert(_ context.Context, p Plant) (Plant, error) {
	if m.err != nil {
		return Plant{}, m.err
	}
	if m.nextID == 0 {
		m.nextID = 1
	}
	p.ID = m.nextID
	m.nextID++
	cp := p
	m.plant = &cp
	return p, nil
}

func (m *memoryStore) Update(_ context.Context, p Plant) (Plant, error) {
	if m.err != nil {
		return Plant{}, m.err
	}
	if m.plant == nil {
		return Plant{}, ErrPlantNotFound
	}
	cp := p
	m.plant = &cp
	return p, nil
}

func TestServicePreservesNotFound(t *testing.T) {
	s := NewService(&memoryStore{})
	_, err := s.GetPlant(context.Background())
	if !errors.Is(err, ErrPlantNotFound) {
		t.Fatalf("err = %v, want ErrPlantNotFound", err)
	}
}

func TestSavePlantInsertsWhenMissing(t *testing.T) {
	store := &memoryStore{}
	s := NewService(store)

	in := PlantFromResponse(Response{
		Epoch:         0,
		Phase:         1,
		PhaseProgress: 0.37,
		Branching:     0.48,
		Density:       0.31,
		Curvature:     0.22,
		Vitality:      0.94,
		Seed:          12345,
	})

	got, err := s.SavePlant(context.Background(), in)
	if err != nil {
		t.Fatalf("SavePlant: %v", err)
	}
	if got.ID == 0 {
		t.Fatal("expected inserted id")
	}
	assertContractJSON(t, got, `{"epoch":0,"phase":1,"phaseProgress":0.37,"branching":0.48,"density":0.31,"curvature":0.22,"vitality":0.94,"seed":12345}`)
}

func TestSavePlantUpdatesExisting(t *testing.T) {
	existing := Plant{ID: 7, Epoch: 0, Phase: 1, Seed: 1}
	s := NewService(&memoryStore{plant: &existing})

	got, err := s.SavePlant(context.Background(), Plant{
		Epoch:         2,
		Phase:         3,
		PhaseProgress: 10,
		Branching:     20,
		Density:       30,
		Curvature:     40,
		Vitality:      50,
		Seed:          99,
	})
	if err != nil {
		t.Fatalf("SavePlant: %v", err)
	}
	if got.ID != 7 {
		t.Fatalf("id = %d, want 7", got.ID)
	}
	if got.Phase != 3 || got.Seed != 99 || got.Branching != 20 {
		t.Fatalf("unexpected plant: %+v", got)
	}
}

func TestResetPlantInsertsInitialWhenMissing(t *testing.T) {
	s := NewService(&memoryStore{})

	got, err := s.ResetPlant(context.Background())
	if err != nil {
		t.Fatalf("ResetPlant: %v", err)
	}

	want := InitialPlant()
	want.ID = got.ID
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestResetPlantRestoresInitialAndKeepsSeed(t *testing.T) {
	existing := Plant{
		ID:            4,
		Epoch:         2,
		Phase:         5,
		PhaseProgress: 37,
		Branching:     48,
		Density:       31,
		Curvature:     22,
		Vitality:      94,
		Seed:          12345,
	}
	s := NewService(&memoryStore{plant: &existing})

	got, err := s.ResetPlant(context.Background())
	if err != nil {
		t.Fatalf("ResetPlant: %v", err)
	}

	want := InitialPlant()
	want.ID = 4
	want.Seed = 12345
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	assertContractJSON(t, got, `{"epoch":0,"phase":1,"phaseProgress":0,"branching":0,"density":0,"curvature":0,"vitality":0,"seed":12345}`)
}

func TestPlantContractRoundTrip(t *testing.T) {
	in := Response{
		Epoch:         0,
		Phase:         1,
		PhaseProgress: 0.37,
		Branching:     0.48,
		Density:       0.31,
		Curvature:     0.22,
		Vitality:      0.94,
		Seed:          12345,
	}
	got := NewResponse(PlantFromResponse(in))
	if got != in {
		t.Fatalf("got %+v, want %+v", got, in)
	}
}

func assertContractJSON(t *testing.T, p Plant, want string) {
	t.Helper()
	raw, err := json.Marshal(NewResponse(p))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != want {
		t.Fatalf("json = %s, want %s", raw, want)
	}
}
