package assignment

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, assignment entity.Assignment) error
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error)
	DeleteByEvent(ctx context.Context, eventID uuid.UUID) error
}
