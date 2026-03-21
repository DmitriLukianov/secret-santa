package user

import (
	"context"
	"secret-santa-backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, user entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, id string, name, email *string) error
	Delete(ctx context.Context, id string) error
}
