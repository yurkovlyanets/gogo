package server

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func newPostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := viper.GetString("postgresql.dsn")
	if dsn == "" {
		return nil, fmt.Errorf("postgresql.dsn is empty")
	}

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse postgres dsn: %w", err)
	}

	cfg.MaxConns = int32(viper.GetInt("postgresql.max_conns"))
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 10
	}
	cfg.MinConns = int32(viper.GetInt("postgresql.min_conns"))
	cfg.MaxConnLifetime = time.Duration(viper.GetInt("postgresql.max_conn_lifetime_sec")) * time.Second
	if cfg.MaxConnLifetime == 0 {
		cfg.MaxConnLifetime = 30 * time.Minute
	}
	cfg.MaxConnIdleTime = time.Duration(viper.GetInt("postgresql.max_conn_idle_sec")) * time.Second
	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = 5 * time.Minute
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping: %w", err)
	}

	return pool, nil
}
