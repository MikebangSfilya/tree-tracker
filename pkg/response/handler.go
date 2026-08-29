package pkg_response

import (
	"encoding/json"
	"log/slog"
	"net/http"
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

func (h *HTTPResponseHandler) WriteError(statusCode int, message string) {
	response := map[string]any{
		"message": message,
		"status":  statusCode,
	}

	h.WriteJSON(response, statusCode)
}
