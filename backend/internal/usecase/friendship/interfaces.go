package friendship

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, f entity.Friendship) (entity.Friendship, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Friendship, error)
	GetByUsers(ctx context.Context, userA, userB uuid.UUID) (*entity.Friendship, error)
	GetFriends(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error)
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
