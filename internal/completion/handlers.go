package completion

import (
	"context"
	"log/slog"
	"net/http"

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

func (h CompletionHandler) Routes(r chi.Router) {
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
		responseHandler.ErrorResponse(err)
		return
	}

	response := make([]CompletionResponse, 0, len(completions))
	for _, completion := range completions {
		response = append(response, CompletionResponse{
			RoutineID:   completion.RoutineID,
			CompletedAt: completion.CompletedAt,
		})
	}

	responseHandler.WriteJSON(response, http.StatusOK)

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

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001") // временное решение, т.к айди придет с контекста
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
