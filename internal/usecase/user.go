package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/yurkovlyanets/gogo/internal/domain"
)

var ErrNotFound = errors.New("not found")

type UserRepo interface {
	Create(ctx context.Context, u domain.User) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (domain.User, error) {
	u, err := domain.NewUser(name, email)
	if err != nil {
		return domain.User{}, err
	}
	return s.repo.Create(ctx, u)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return s.repo.GetByID(ctx, id)
}
