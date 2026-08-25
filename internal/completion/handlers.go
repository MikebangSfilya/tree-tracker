package completion

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	completion_errors "github.com/MikebangSfilya/tree-tracker/pkg/errors"
	pkg_response "github.com/MikebangSfilya/tree-tracker/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CompletionHandler struct {
	completionService CompletionService
	logger            *slog.Logger
}

type CompletionService interface {
	RoutineComplete(ctx context.Context, input CompleteInput) (Completion, error)
	GetCompletions(ctx context.Context, userID uuid.UUID) ([]Completion, error)
}

func NewCompletionHandler(completionService CompletionService) *CompletionHandler {
	return &CompletionHandler{
		completionService: completionService,
		logger:            slog.Default(),
	}
}

func (h *CompletionHandler) Routes(r chi.Router) {
	r.Post("/api/routines/{id}/complete", h.RoutineComplete)
	r.Get("/api/completions", h.GetCompletions)
}

func (h *CompletionHandler) GetCompletions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	responseHandler := pkg_response.NewHTTPResponseHandler(h.logger, w)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	var completions []Completion

	completions, err := h.completionService.GetCompletions(ctx, userID)
	if err != nil {
		h.logger.Error("failed to get completions", "error", err)
		h.errorResponse(responseHandler, err)
		return
	}

	response := make([]CompletionResponse, 0, len(completions))
	for _, completion := range completions {
		response = append(response, CompletionResponse{
			RoutineID:   completion.RoutineId,
			CompletedAt: completion.CompletedAt,
		})
	}

	responseHandler.WriteJSON(response, http.StatusOK)

}

func (h *CompletionHandler) RoutineComplete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	responseHandler := pkg_response.NewHTTPResponseHandler(h.logger, w)
	routineId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error("failed to parse routineID", slog.String("id", idStr), slog.String("error", err.Error()))
		responseHandler.WriteError(http.StatusBadRequest, "invalid routine id, must be an integer")
		return
	}

	userId := uuid.MustParse("00000000-0000-0000-0000-000000000001") // временное решение, т.к айди придет с контекста
	input := NewCompleteInput(userId, routineId)

	completion, err := h.completionService.RoutineComplete(ctx, input)
	if err != nil {
		h.logger.Error("failed to complete routine")
		h.errorResponse(responseHandler, err)
		return
	}

	responseHandler.WriteJSON(
		CompletionResponse{
			RoutineID:   completion.RoutineId,
			CompletedAt: completion.CompletedAt,
		},
		http.StatusOK)
}

func (h *CompletionHandler) errorResponse(resp *pkg_response.HTTPResponseHandler, err error) {
	var (
		statusCode int
		msg        string
	)

	switch {
	case errors.Is(err, completion_errors.ErrAlreadyCompleted):
		statusCode = http.StatusConflict
		msg = "routine is already completed"
	case errors.Is(err, completion_errors.ErrAlreadyCompletedToday):
		statusCode = http.StatusConflict
		msg = "routine had already been completed today"
	default:
		statusCode = http.StatusInternalServerError
		msg = "failed to complete routine"
	}

	resp.WriteError(statusCode, msg)
}
