package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yurkovlyanets/gogo/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

// Create сохраняет пользователя и возвращает сущность с заполненными полями из БД (например created_at).
func (r *UserRepo) Create(ctx context.Context, u domain.User) (domain.User, error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	// Если created_at у тебя заполняется в БД через DEFAULT now(),
	// можно не передавать его и вернуть через RETURNING.
	const q = `
		INSERT INTO users (id, name, email)
		VALUES ($1, $2, $3)
		RETURNING created_at
	`

	var createdAt time.Time
	err := r.pool.QueryRow(ctx, q, u.ID, u.Name, u.Email).Scan(&createdAt)
	if err != nil {
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	u.CreatedAt = createdAt
	return u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	const q = `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	var u domain.User
	err := r.pool.QueryRow(ctx, q, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("select user by id: %w", err)
	}

	return u, nil
}
