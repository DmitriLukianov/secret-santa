package repository

import (
	"context"
	"secret-santa-backend/internal/domain"
)

type ParticipantRepository interface {
	AddParticipant(ctx context.Context, participant domain.Participant) error
	GetParticipantsByEvent(ctx context.Context, eventID string) ([]domain.Participant, error)
	DeleteParticipant(ctx context.Context, eventID, userID string) error
}
