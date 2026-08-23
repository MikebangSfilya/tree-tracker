package pkg_response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	completion_errors "github.com/MikebangSfilya/tree-tracker/pkg/completion errors"
)

type HTTPResponseHandler struct {
	logger *slog.Logger
	w      http.ResponseWriter
}

func NewHTTPResponseHandler(logger *slog.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		logger: logger,
		w:      w,
	}
}

func (h *HTTPResponseHandler) WriteJSON(data any, statusCode int) {
	h.w.Header().Set("Content-Type", "application/json")
	h.w.WriteHeader(statusCode)

	if err := json.NewEncoder(h.w).Encode(data); err != nil {
		h.logger.Error("failed to encode and write json response", "err", err)
	}
}

func (h *HTTPResponseHandler) ErrorResponse(err error) {
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

	h.writeError(statusCode, msg, err)
}

func (h *HTTPResponseHandler) writeError(code int, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	response := map[string]string{
		"error":   errMsg,
		"message": message,
	}

	h.WriteJSON(response, code)
}
