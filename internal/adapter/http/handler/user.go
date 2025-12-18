package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	httpx "github.com/yurkovlyanets/gogo/internal/adapter/http/response"
	"github.com/yurkovlyanets/gogo/internal/usecase"
)

type UserHandler struct {
	svc *usecase.UserService
}

func NewUserHandler(svc *usecase.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteBadRequest(w, "invalid json")
		return
	}

	u, err := h.svc.CreateUser(r.Context(), req.Name, req.Email)
	if err != nil {
		httpx.WriteBadRequest(w, err.Error())
		return
	}

	httpx.WriteCreated(w, userResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(timeRFC3339),
	})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		httpx.WriteBadRequest(w, "invalid id")
		return
	}

	u, err := h.svc.GetUserByID(r.Context(), id)
	if errors.Is(err, usecase.ErrNotFound) {
		httpx.WriteNotFound(w, "user not found")
		return
	}
	if err != nil {
		httpx.WriteInternal(w, "internal error")
		return
	}

	httpx.WriteOKJSON(w, userResponse{
		ID:        u.ID.String(),
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Format(timeRFC3339),
	})
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"
