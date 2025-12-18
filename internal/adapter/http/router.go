package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yurkovlyanets/gogo/internal/adapter/http/handler"
)

type RouterDeps struct {
	UserHandler   *handler.UserHandler
	HealthHandler *handler.HealthHandler
}

func NewRouter(deps RouterDeps) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", deps.HealthHandler.Health)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", deps.UserHandler.Create)
			r.Get("/{id}", deps.UserHandler.GetByID)
		})
	})

	return r
}
