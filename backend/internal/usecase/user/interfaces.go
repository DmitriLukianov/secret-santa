package user

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error)
	Update(ctx context.Context, id uuid.UUID, name, email *string) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
