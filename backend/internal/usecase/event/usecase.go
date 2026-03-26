package event

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, log: log}
}

// Create — создаёт событие
func (uc *UseCase) Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error) {
	now := time.Now()
	if uc.log != nil {
		uc.log.Info("create event started",
			slog.String("organizer_id", organizerID.String()),
			slog.String("title", input.Title),
		)
	}

	// DTO имеет StartDate/EndDate как time.Time (не указатели), DrawDate — *time.Time
	startDate := now
	if !input.StartDate.IsZero() {
		startDate = input.StartDate
	}

	drawDate := now
	if input.DrawDate != nil {
		drawDate = *input.DrawDate
	}

	endDate := now
	if !input.EndDate.IsZero() {
		endDate = input.EndDate
	}

	event := entity.NewEvent(
		input.Title,
		organizerID,
		input.Description,
		input.Rules,
		input.Recommendations,
		startDate,
		drawDate,
		endDate,
		input.MaxParticipants,
	)

	if err := uc.repo.Create(ctx, event); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create event",
				slog.String("organizer_id", organizerID.String()),
				slog.String("title", input.Title),
				slog.String("error", err.Error()),
			)
		}
		return entity.Event{}, fmt.Errorf("failed to create event: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("event created successfully",
			slog.String("event_id", event.ID.String()),
			slog.String("organizer_id", organizerID.String()),
		)
	}

	return event, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("event id is required")
	}
	if uc.log != nil {
		uc.log.Info("get event by id started", slog.String("event_id", id.String()))
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Event, error) {
	if uc.log != nil {
		uc.log.Info("get all events started")
	}
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error {
	if id == uuid.Nil {
		return fmt.Errorf("event id is required")
	}
	if uc.log != nil {
		uc.log.Info("update event started",
			slog.String("event_id", id.String()),
			slog.Any("input", input),
		)
	}
	if err := uc.repo.Update(ctx, id, input); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update event",
				slog.String("event_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	if uc.log != nil {
		uc.log.Info("event updated successfully", slog.String("event_id", id.String()))
	}
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("event id is required")
	}
	if uc.log != nil {
		uc.log.Info("delete event started", slog.String("event_id", id.String()))
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete event",
				slog.String("event_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	if uc.log != nil {
		uc.log.Info("event deleted successfully", slog.String("event_id", id.String()))
	}
	return nil
}

// Finish — завершает событие (только организатор)
func (uc *UseCase) Finish(ctx context.Context, id, userID uuid.UUID) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return fmt.Errorf("event id and user id are required")
	}
	if uc.log != nil {
		uc.log.Info("finish event started",
			slog.String("event_id", id.String()),
			slog.String("user_id", userID.String()),
		)
	}

	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get event for finish",
				slog.String("event_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to get event: %w", err)
	}

	// Проверка организатора
	if event.OrganizerID != userID {
		if uc.log != nil {
			uc.log.Warn("finish denied: not organizer",
				slog.String("event_id", id.String()),
				slog.String("user_id", userID.String()),
				slog.String("organizer_id", event.OrganizerID.String()),
			)
		}
		return fmt.Errorf("only the event organizer can finish the event")
	}

	if !event.CanBeFinished() {
		if uc.log != nil {
			uc.log.Warn("finish denied: invalid event status",
				slog.String("event_id", id.String()),
				slog.String("status", string(event.Status)),
			)
		}
		return fmt.Errorf("event cannot be finished in current status: %s", event.Status)
	}

	event.MarkAsFinished()

	if err := uc.repo.Update(ctx, id, dto.UpdateEventInput{Status: stringPtr(string(entity.EventStatusFinished))}); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to mark event as finished",
				slog.String("event_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return err
	}

	if uc.log != nil {
		uc.log.Info("event finished successfully",
			slog.String("event_id", id.String()),
			slog.String("user_id", userID.String()),
		)
	}
	return nil
}

// stringPtr — маленькая вспомогательная функция
func stringPtr(s string) *string {
	return &s
}
