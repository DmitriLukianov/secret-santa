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

	// FIXED: новый метод — вся жеребьёвка в одной атомарной транзакции
	TransactionalDraw(ctx context.Context, eventID uuid.UUID, assignments []entity.Assignment, newStatus entity.EventStatus) error
}

type EventRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.EventStatus) error
}

type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}
