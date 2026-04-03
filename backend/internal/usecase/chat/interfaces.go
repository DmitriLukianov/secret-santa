package chat

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

// Repository — локальный интерфейс репозитория для чата
// (полностью соответствует ChatRepository из contracts.go)
type Repository interface {
	// CreateMessage теперь возвращает полностью заполненную сущность из БД
	CreateMessage(ctx context.Context, msg entity.Message) (entity.Message, error)
	GetMessagesByPair(ctx context.Context, eventID, user1ID, user2ID uuid.UUID) ([]entity.Message, error)
}

// ParticipantRepository и AssignmentRepository — зависимости usecase
type ParticipantRepository interface {
	// Пока не используется в chat, но оставляем для единообразия
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}

type AssignmentRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error)
}
