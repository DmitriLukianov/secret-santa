package event

import (
	"context"
	"log/slog"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{
		repo: repo,
		log:  log,
	}
}

func (uc *UseCase) Create(ctx context.Context, input dto.CreateEventInput, organizerID string) (entity.Event, error) {
	// Здесь можно добавить валидацию дат и бизнес-правила позже

	event := entity.NewEvent(
		input.Name,
		uuid.MustParse(organizerID),
		input.Description,
		input.Rules,
		input.Recommendations,
		input.StartDate,
		input.DrawDate,
		input.EndDate,
		input.MaxParticipants,
	)

	created, err := uc.repo.Create(ctx, event)
	if err != nil {
		uc.log.Error("failed to create event", "error", err)
		return entity.Event{}, err
	}

	uc.log.Info("event created", "event_id", created.ID, "name", created.Name)
	return created, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id string) (entity.Event, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) List(ctx context.Context) ([]entity.Event, error) {
	return uc.repo.List(ctx)
}

func (uc *UseCase) Update(ctx context.Context, id string, input dto.UpdateEventInput) (entity.Event, error) {
	// Заглушка — позже реализуем полное обновление
	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Event{}, err
	}

	// Применяем изменения из input (partial update)
	if input.Name != nil {
		event.Name = *input.Name
	}
	if input.Description != nil {
		event.Description = input.Description
	}
	if input.Rules != nil {
		event.Rules = input.Rules
	}
	if input.Recommendations != nil {
		event.Recommendations = input.Recommendations
	}
	if input.StartDate != nil {
		event.StartDate = input.StartDate
	}
	if input.DrawDate != nil {
		event.DrawDate = input.DrawDate
	}
	if input.EndDate != nil {
		event.EndDate = input.EndDate
	}
	if input.Status != nil {
		event.Status = *input.Status
	}
	if input.MaxParticipants != nil {
		event.MaxParticipants = *input.MaxParticipants
	}

	return uc.repo.Update(ctx, event)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
