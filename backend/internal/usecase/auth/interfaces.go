package auth

import (
	"context"
	"secret-santa-backend/internal/entity"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID string) (*entity.User, error)
	Create(ctx context.Context, user entity.User) error
}
