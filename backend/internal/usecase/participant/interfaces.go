package participant

import (
	"context"
	"secret-santa-backend/internal/entity"
)

type Repository interface {
	Add(ctx context.Context, participant entity.Participant) error
	GetByEvent(ctx context.Context, eventID string) ([]entity.Participant, error)
	Delete(ctx context.Context, id string) error
}
