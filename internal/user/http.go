package user

import (
	"encoding/json"
	"errors"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

func (h Handler) writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		h.logger.Error("failed to write json response", "error", err)
	}
}

func (h Handler) writeError(w http.ResponseWriter, err error) {
	status, msg := httpStatus(err)
	if status >= http.StatusInternalServerError {
		h.logger.Error("user request failed", "error", err)
	}
	h.writeJSON(w, status, errorBody{Error: msg})
}

func (h Handler) writeBadRequest(w http.ResponseWriter, msg string) {
	h.writeJSON(w, http.StatusBadRequest, errorBody{Error: msg})
}

func httpStatus(err error) (int, string) {
	switch {
	case errors.Is(err, ErrUserAlreadyExists):
		return http.StatusConflict, "user already exists"
	case errors.Is(err, ErrUserNotFound):
		return http.StatusNotFound, "user not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}
