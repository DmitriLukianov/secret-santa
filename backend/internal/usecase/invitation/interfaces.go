package invitation

import (
	"context"

	"secret-santa-backend/internal/entity"
)

// Repository — интерфейс репозитория для приглашений (многоразовые ссылки)
type Repository interface {
	Create(ctx context.Context, invitation entity.Invitation) error
	GetByToken(ctx context.Context, token string) (*entity.Invitation, error)
}
