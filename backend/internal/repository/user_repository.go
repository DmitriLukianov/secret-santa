package repository

import (
	"context"
	"secret-santa-backend/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	UpdateUser(ctx context.Context, id string, name, email *string) error
	Delete(ctx context.Context, id string) error
}
