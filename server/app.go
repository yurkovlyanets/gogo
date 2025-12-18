package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/yurkovlyanets/gogo/config"
	httpx "github.com/yurkovlyanets/gogo/internal/adapter/http"
	"github.com/yurkovlyanets/gogo/internal/adapter/http/handler"
	"github.com/yurkovlyanets/gogo/internal/infra/repo"
	"github.com/yurkovlyanets/gogo/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server
	pool       *pgxpool.Pool
}

func New() (*App, error) {
	if err := config.Init(); err != nil {
		return nil, fmt.Errorf("config init: %w", err)
	}

	port := viper.GetInt("port")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := newPostgresPool(ctx)
	if err != nil {
		return nil, err
	}

	userRepo := repo.NewUserRepo(pool)
	userUC := usecase.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userUC)
	healthHandler := handler.NewHealthHandler()

	handler := httpx.NewRouter(httpx.RouterDeps{userHandler, healthHandler})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		httpServer: srv,
		pool:       pool,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("http listening on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		log.Printf("shutdown: %v", ctx.Err())

		_ = a.httpServer.Shutdown(shutdownCtx)
		a.pool.Close()
		return nil

	case err := <-errCh:
		// если сервер упал — закрываем pool
		a.pool.Close()
		return err
	}
}
