package plant

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeGetter struct {
	plant Plant
	err   error
}

func (f fakeGetter) GetPlant(context.Context) (Plant, error) {
	return f.plant, f.err
}

func TestGetPlantOK(t *testing.T) {
	h := NewHandler(fakeGetter{plant: Plant{
		Epoch:         0,
		Phase:         1,
		PhaseProgress: 37,
		Branching:     48,
		Density:       31,
		Curvature:     22,
		Vitality:      94,
		Seed:          12345,
	}}, nil)

	rr := serve(t, h, http.MethodGet, "/api/plant")
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}

	got := strings.TrimSpace(rr.Body.String())
	want := `{"epoch":0,"phase":1,"phaseProgress":0.37,"branching":0.48,"density":0.31,"curvature":0.22,"vitality":0.94,"seed":12345}`
	if got != want {
		t.Fatalf("body = %s, want %s", got, want)
	}
}

func TestGetPlantNotFound(t *testing.T) {
	h := NewHandler(fakeGetter{err: fmt.Errorf("get plant: %w", ErrPlantNotFound)}, nil)

	rr := serve(t, h, http.MethodGet, "/api/plant")
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotFound)
	}

	assertErrorBody(t, rr, "plant not found")
}

func TestGetPlantInternalErrorHidesCause(t *testing.T) {
	h := NewHandler(fakeGetter{err: errors.New("SQLSTATE 42P01")}, nil)

	rr := serve(t, h, http.MethodGet, "/api/plant")
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	assertErrorBody(t, rr, "internal server error")
}

func TestNewResponseScalesRatios(t *testing.T) {
	got := NewResponse(Plant{
		PhaseProgress: 37,
		Branching:     48,
		Density:       31,
		Curvature:     22,
		Vitality:      94,
		Seed:          7,
	})

	raw, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"epoch":0,"phase":0,"phaseProgress":0.37,"branching":0.48,"density":0.31,"curvature":0.22,"vitality":0.94,"seed":7}`
	if string(raw) != want {
		t.Fatalf("json = %s, want %s", raw, want)
	}
}

func serve(t *testing.T, h *Handler, method, path string) *httptest.ResponseRecorder {
	t.Helper()

	r := chi.NewRouter()
	h.Routes(r)

	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func assertErrorBody(t *testing.T, rr *httptest.ResponseRecorder, want string) {
	t.Helper()

	var body errorBody
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Error != want {
		t.Fatalf("error = %q, want %q", body.Error, want)
	}
}
