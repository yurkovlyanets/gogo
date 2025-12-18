package domain

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
}

func NewUser(name, email string) (User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if name == "" {
		return User{}, errors.New("name is required")
	}
	if email == "" {
		return User{}, errors.New("email is required")
	}

	return User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
		// CreatedAt заполнится в БД, но пусть будет ок и в памяти
		CreatedAt: time.Now().UTC(),
	}, nil
}
