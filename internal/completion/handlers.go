package completion

import (
	"context"
	"log/slog"
	"net/http"

	pkg_response "github.com/MikebangSfilya/tree-tracker/pkg/response"
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
	responseHandler := pkg_response.NewHTTPResponseHandler(h.logger, w)
	routineID, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Error("failed to parse routineID", slog.String("id", idStr), slog.String("error", err.Error()))
		responseHandler.ErrorResponse(err)
		return
	}

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	input := NewCompleteInput(routineID, userID)

	completion, err := h.completionService.RoutineComplete(ctx, input)
	if err != nil {
		h.logger.Error("failed to complete routine")
		responseHandler.ErrorResponse(err)
		return
	}

	responseHandler.WriteJSON(
		CompletionResponse{
			RoutineID:   completion.RoutineID,
			CompletedAt: completion.CompletedAt,
		},
		http.StatusOK)
}
