package participant

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Add(ctx context.Context, input dto.AddParticipantInput) error {
	if input.EventID == "" {
		return fmt.Errorf("event_id is required")
	}
	if input.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	p := entity.Participant{
		ID:      uuid.NewString(),
		EventID: input.EventID,
		UserID:  input.UserID,
	}

	return uc.repo.Add(ctx, p)
}
func (uc *UseCase) GetByEvent(ctx context.Context, eventID string) ([]entity.Participant, error) {
	if eventID == "" {
		return nil, fmt.Errorf("event_id is required")
	}

	return uc.repo.GetByEvent(ctx, eventID)
}
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return uc.repo.Delete(ctx, id)
}
