package event

import (
	"context"
	"time"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, event entity.Event) (entity.Event, error)

	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status definitions.EventStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetEventsForUser(ctx context.Context, userID uuid.UUID) ([]entity.Event, error)
}

type Scheduler interface {
	Schedule(eventID uuid.UUID, drawAt time.Time)
	Cancel(eventID uuid.UUID)
}
