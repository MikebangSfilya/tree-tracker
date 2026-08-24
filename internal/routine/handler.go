package routine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	routines, err := h.service.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, NewResponses(routines))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var request CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	routine, err := h.service.Create(r.Context(), request)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, NewResponse(routine))
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid routine id"))
		return
	}

	var request UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	routine, err := h.service.Update(r.Context(), id, request)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, NewResponse(routine))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid routine id"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getID(r *http.Request) (int64, error) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) == 0 {
		return 0, errors.New("missing id")
	}

	idPart := parts[len(parts)-1]

	return strconv.ParseInt(idPart, 10, 64)
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrRoutineNotFound):
		writeError(w, http.StatusNotFound, err)

	case errors.Is(err, ErrInvalidCategory),
		errors.Is(err, ErrInvalidType),
		errors.Is(err, ErrInvalidWeight),
		errors.Is(err, ErrInvalidCoefficient):
		writeError(w, http.StatusBadRequest, err)

	default:
		writeError(w, http.StatusInternalServerError, errors.New("internal server error"))
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{
		"error": err.Error(),
	})
}
