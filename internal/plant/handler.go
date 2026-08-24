package plant

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Getter interface {
	GetPlant(ctx context.Context) (Plant, error)
}

type Handler struct {
	service Getter
	logger  *slog.Logger
}

func NewHandler(service Getter, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/api/plant", h.getPlant)
}

func (h *Handler) getPlant(w http.ResponseWriter, r *http.Request) {
	p, err := h.service.GetPlant(r.Context())
	if err != nil {
		h.writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, NewResponse(p))
}

type errorBody struct {
	Error string `json:"error"`
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	status, msg := httpStatus(err)
	if status >= http.StatusInternalServerError {
		h.logger.Error("plant request failed", "error", err)
	}
	writeJSON(w, status, errorBody{Error: msg})
}

func httpStatus(err error) (int, string) {
	switch {
	case errors.Is(err, ErrPlantNotFound):
		return http.StatusNotFound, "plant not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("failed to write json response", "error", err)
	}
}
