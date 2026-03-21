package event

import (
	"context"
	"fmt"
	"time"

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

// CREATE EVENT
func (uc *UseCase) Create(ctx context.Context, input dto.CreateEventInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.OrganizerID == "" {
		return fmt.Errorf("organizer_id is required")
	}

	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		return fmt.Errorf("invalid start_date")
	}

	drawDate, err := time.Parse(time.RFC3339, input.DrawDate)
	if err != nil {
		return fmt.Errorf("invalid draw_date")
	}

	endDate, err := time.Parse(time.RFC3339, input.EndDate)
	if err != nil {
		return fmt.Errorf("invalid end_date")
	}

	event := entity.Event{
		ID:          uuid.NewString(),
		Name:        input.Name,
		Description: input.Description,
		OrganizerID: input.OrganizerID,
		StartDate:   startDate,
		DrawDate:    drawDate,
		EndDate:     endDate,
	}

	return uc.repo.Create(ctx, event)
}

// GET EVENT BY ID
func (uc *UseCase) Get(ctx context.Context, id string) (*entity.Event, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	return uc.repo.GetByID(ctx, id)
}

// GET ALL EVENTS
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Event, error) {
	return uc.repo.GetAll(ctx)
}

// UPDATE EVENT
func (uc *UseCase) Update(ctx context.Context, id string, input dto.UpdateEventInput) error {
	if input.Name == nil && input.Description == nil {
		return fmt.Errorf("nothing to update")
	}
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return uc.repo.Update(ctx, id, input.Name, input.Description)
}

// DELETE EVENT
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return uc.repo.Delete(ctx, id)
}
