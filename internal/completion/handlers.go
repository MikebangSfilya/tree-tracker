package completion

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type CompletionHandler struct {
	completionService CompletionService
	logger            *slog.Logger
}

type CompletionService interface {
	RoutineComplete(ctx context.Context, input CompleteInput) (Completion, error)
}

func NewCompletionHandler(completionService CompletionService) *CompletionHandler {
	return &CompletionHandler{
		completionService: completionService,
		logger:            slog.Default(),
	}
}

func (h *CompletionHandler) RoutineComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := r.PathValue("id")

	routineID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to parse routineID", slog.String("id", idStr), slog.String("error", err.Error()))
		writeError(w, http.StatusBadRequest, "failed to parse id", err)
		return
	}

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	input := NewCompleteInput(routineID, userID)

	completion, err := h.completionService.RoutineComplete(ctx, input)
	if err != nil {
		h.logger.Error("failed to complete routine")
		errorResponse(w, err)
		return
	}

	writeJSON(w,
		CompletionResponse{
			RoutineID:   completion.RoutineID,
			CompletedAt: completion.CompletedAt,
		},
		http.StatusOK)

}

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)

}

func errorResponse(w http.ResponseWriter, err error) {
	var (
		statusCode int
		msg        string
	)

	switch {
	case errors.Is(err, ErrAlreadyCompleted):
		statusCode = http.StatusConflict
		msg = "routine is already completed"
	case errors.Is(err, ErrAlreadyCompletedToday):
		statusCode = http.StatusConflict
		msg = "routine had already been completed today"
	default:
		statusCode = http.StatusInternalServerError
		msg = "failed to complete routine"
	}

	writeError(w, statusCode, msg, err)
}

func writeError(w http.ResponseWriter, code int, message string, err error) {
	response := map[string]string{
		"error":   err.Error(),
		"message": message,
	}

	writeJSON(w, response, code)
}
