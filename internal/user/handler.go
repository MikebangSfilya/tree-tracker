package user

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServiceMethods interface {
	Register(ctx context.Context, username, email, password string) (id string, err error)
	GetUser(ctx context.Context, email string) (user *User, err error)
	DeleteUser(ctx context.Context, email string) (err error)
}

type Handler struct {
	logger *slog.Logger

	srv ServiceMethods
}

func NewHandler(logger *slog.Logger, srv ServiceMethods) *Handler {
	return &Handler{
		logger: logger,
		srv:    srv,
	}
}

func (h Handler) Routes(r chi.Router) {
	r.Post("/api/user", h.register)
	r.Get("/api/user", h.getUser)
	r.Delete("/api/user", h.deleteUser)
}

func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var userData CreateReq
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		h.writeBadRequest(w, "invalid request body")
		return
	}

	id, err := h.srv.Register(
		r.Context(),
		userData.Username,
		userData.Email,
		userData.Password,
	)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, CreateResp{ID: id})
}

func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	var emailData GetAndDeleteReq
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		h.writeBadRequest(w, "invalid request body")
		return
	}

	user, err := h.srv.GetUser(r.Context(), emailData.Email)
	if err != nil {
		h.writeError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, GetResp{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var emailData GetAndDeleteReq
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		h.writeBadRequest(w, "invalid request body")
		return
	}

	err = h.srv.DeleteUser(r.Context(), emailData.Email)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
