package event

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

// Create — создаёт событие
func (uc *UseCase) Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error) {
	event := entity.NewEvent(
		input.Title,
		organizerID,
		input.Description,
		input.Rules,
		input.Recommendations,
		&input.StartDate,
		input.DrawDate,
		&input.EndDate,
		input.MaxParticipants,
	)

	if err := uc.repo.Create(ctx, event); err != nil {
		return entity.Event{}, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("event id is required")
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Event, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error {
	if id == uuid.Nil {
		return fmt.Errorf("event id is required")
	}
	return uc.repo.Update(ctx, id, input)
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("event id is required")
	}
	return uc.repo.Delete(ctx, id)
}
