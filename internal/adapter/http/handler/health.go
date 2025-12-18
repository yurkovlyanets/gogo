package handler

import (
	"net/http"

	httpx "github.com/yurkovlyanets/gogo/internal/adapter/http/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	httpx.WriteOK(w)
}
