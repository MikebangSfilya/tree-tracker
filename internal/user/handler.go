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


func NewHandler(srv ServiceMethods) *Handler {
	return &Handler{srv: srv}
}


func (h Handler) Routes(r chi.Router) {
	// assuming basic route is `/api/`
	r.Post("/user", h.register)
	r.Get("/user", h.getUser)
	r.Delete("/user", h.deleteUser)
}


func (h Handler) register(w http.ResponseWriter, r *http.Request) {
	var userData CreateReq
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.srv.Register(
		r.Context(),
		userData.Username,
		userData.Email,
		userData.Password,
	)
	if err != nil {
		if err == ErrUserAlreadyExists {
			http.Error(w, ErrUserAlreadyExists.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var idResp = CreateResp{ID: id}
	err = json.NewEncoder(w).Encode(idResp)
	if err != nil {
		h.logger.Error("failed to encode",
			"id", idResp.ID,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (h Handler) getUser(w http.ResponseWriter, r *http.Request) {
	var emailData GetAndDeleteReq
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.srv.GetUser(r.Context(), emailData.Email)
	if err != nil {
		if err == ErrUserNotFound {
			http.Error(w, ErrUserNotFound.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var userResp = GetResp{
		ID: 	   user.ID.String(),
		Username:  user.Username,
		Email: 	   user.Email,
		CreatedAt: user.CreatedAt.String(),
	}
	err = json.NewEncoder(w).Encode(userResp)
	if err != nil {
		h.logger.Error("failed to encode",
			"id", userResp.ID,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (h Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	var emailData GetAndDeleteReq
	err := json.NewDecoder(r.Body).Decode(&emailData)
	if err != nil {
		h.logger.Error("failed to decode", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.srv.DeleteUser(r.Context(), emailData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var deleteResp = DeleteResp{Deleted: true}
	err = json.NewEncoder(w).Encode(deleteResp)
	if err != nil {
		h.logger.Error("failed to encode",
			"email", emailData.Email,
			"error", err,
		)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
