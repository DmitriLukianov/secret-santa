package assignment

import (
	"context"

	"secret-santa-backend/internal/entity"
)

type AssignmentRepository interface {
	CreateMany(ctx context.Context, assignments []entity.Assignment) error
	GetByEvent(ctx context.Context, eventID string) ([]entity.Assignment, error)
}

type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID string) ([]entity.Participant, error)
}
