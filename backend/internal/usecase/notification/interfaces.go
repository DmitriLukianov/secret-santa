package notification

import (
	"context"

	"secret-santa-backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, n entity.Notification) (entity.Notification, error)
}
