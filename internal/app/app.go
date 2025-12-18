package app

import (
	"context"
	"net/http"

	httpx "github.com/yurkovlyanets/gogo/internal/adapter/http"
	"github.com/yurkovlyanets/gogo/internal/adapter/http/handler"
	"github.com/yurkovlyanets/gogo/internal/infra/repo"
	"github.com/yurkovlyanets/gogo/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHTTPHandler(ctx context.Context, pool *pgxpool.Pool) http.Handler {
	userRepo := repo.NewUserPostgresRepo(pool)
	userSvc := usecase.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userSvc)
	healthHandler := handler.NewHealthHandler()

	return httpx.NewRouter(httpx.RouterDeps{
		UserHandler:   userHandler,
		HealthHandler: healthHandler,
	})
}
