package participant

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Create — добавляет участника в событие
func (uc *UseCase) Create(ctx context.Context, eventID, userID uuid.UUID, role string) (entity.Participant, error) {
	participant := entity.NewParticipant(eventID, userID, role)

	if err := uc.repo.Create(ctx, participant); err != nil {
		return entity.Participant{}, fmt.Errorf("failed to create participant: %w", err)
	}

	return participant, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("participant id is required")
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error) {
	if eventID == uuid.Nil {
		return nil, fmt.Errorf("event id is required")
	}
	return uc.repo.GetByEvent(ctx, eventID)
}

// MarkGiftSent — отметка, что участник отправил подарок
func (uc *UseCase) MarkGiftSent(ctx context.Context, participantID uuid.UUID) error {
	if participantID == uuid.Nil {
		return fmt.Errorf("participant id is required")
	}
	return uc.repo.UpdateGiftSent(ctx, participantID, true)
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("participant id is required")
	}
	return uc.repo.Delete(ctx, id)
}
