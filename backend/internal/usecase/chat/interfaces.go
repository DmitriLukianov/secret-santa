package chat

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	CreateMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetMessagesByPair(ctx context.Context, eventID, user1ID, user2ID uuid.UUID) ([]entity.Message, error)
}

type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}

type AssignmentRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error)
}
