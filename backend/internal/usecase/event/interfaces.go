package event

import (
	"context"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, event entity.Event) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}
