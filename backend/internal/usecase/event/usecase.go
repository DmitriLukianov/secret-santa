package event

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
	participant "secret-santa-backend/internal/usecase/participant"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/google/uuid"
)

type UseCase struct {
	repo            Repository
	participantRepo participant.Repository
	log             *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, participantRepo participant.Repository, log *slog.Logger) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		log:             log,
	}
}

func (uc *UseCase) Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error) {
	if uc.log != nil {
		uc.log.Info("create event started",
			slog.String("organizer_id", organizerID.String()),
			slog.String("title", input.Title),
		)
	}

	if organizerID == uuid.Nil {
		return entity.Event{}, definitions.ErrInvalidUserInput
	}

	event := entity.NewEvent(
		input.Title,
		organizerID,
		input.Description,
		input.Rules,
		input.Recommendations,
		input.StartDate,
		input.EndDate,
		input.DrawDate,
		input.MaxParticipants,
	)

	createdEvent, err := uc.repo.Create(ctx, event)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create event", slog.String("error", err.Error()))
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.Event{}, definitions.ErrConflict
		}
		return entity.Event{}, fmt.Errorf("failed to create event: %w", err)
	}

	if input.WantParticipate {
		organizerParticipant := entity.NewParticipant(createdEvent.ID, organizerID, definitions.ParticipantRoleOrganizer)
		if _, err = uc.participantRepo.Create(ctx, organizerParticipant); err != nil {
			if uc.log != nil {
				uc.log.Error("failed to create organizer participant", slog.String("error", err.Error()))
			}
			return entity.Event{}, fmt.Errorf("failed to create organizer participant: %w", err)
		}
	}

	if uc.log != nil {
		uc.log.Info("event created successfully",
			slog.String("event_id", createdEvent.ID.String()),
		)
	}
	return createdEvent, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	if id == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	event, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, definitions.ErrEventNotFound
	}
	return event, nil
}

func (uc *UseCase) GetAll(ctx context.Context) ([]entity.Event, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) Update(ctx context.Context, id, userID uuid.UUID, input dto.UpdateEventInput) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	eventPtr, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return definitions.ErrEventNotFound
	}

	if eventPtr.OrganizerID != userID {
		return definitions.ErrNotOrganizer
	}

	if uc.log != nil {
		uc.log.Info("update event started", slog.String("event_id", id.String()))
	}

	if err := uc.repo.Update(ctx, id, input); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update event", slog.String("error", err.Error()))
		}
		return err
	}

	if uc.log != nil {
		uc.log.Info("event updated successfully", slog.String("event_id", id.String()))
	}
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id, userID uuid.UUID) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	eventPtr, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return definitions.ErrEventNotFound
	}

	if eventPtr.OrganizerID != userID {
		return definitions.ErrNotOrganizer
	}

	if uc.log != nil {
		uc.log.Info("delete event started", slog.String("event_id", id.String()))
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete event", slog.String("error", err.Error()))
		}
		return err
	}

	if uc.log != nil {
		uc.log.Info("event deleted successfully", slog.String("event_id", id.String()))
	}
	return nil
}

func (uc *UseCase) Activate(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusActive)
}

func (uc *UseCase) Finish(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusFinished)
}

func (uc *UseCase) StartDrawing(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusDrawingPending)
}

func (uc *UseCase) OpenInvitation(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusInvitationOpen)
}

func (uc *UseCase) CloseRegistration(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusRegistrationClosed)
}

func (uc *UseCase) Cancel(ctx context.Context, id, userID uuid.UUID) error {
	return uc.changeStatus(ctx, id, userID, definitions.EventStatusCancelled)
}

func (uc *UseCase) changeStatus(ctx context.Context, id, userID uuid.UUID, newStatus definitions.EventStatus) error {
	if id == uuid.Nil || userID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	eventPtr, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return definitions.ErrEventNotFound
	}

	if eventPtr.OrganizerID != userID {
		return definitions.ErrNotOrganizer
	}

	if eventPtr.Status == newStatus {
		return fmt.Errorf("%w: status already %s", definitions.ErrInvalidEventState, newStatus)
	}

	oldStatus := eventPtr.Status
	if err := eventPtr.TransitionTo(newStatus); err != nil {
		return err
	}

	if err := uc.repo.UpdateStatus(ctx, id, eventPtr.Status); err != nil {
		return err
	}

	if uc.log != nil {
		uc.log.Info("event status changed",
			slog.String("event_id", id.String()),
			slog.String("old_status", string(oldStatus)),
			slog.String("new_status", string(newStatus)),
		)
	}
	return nil
}

func (uc *UseCase) GetMyEvents(ctx context.Context, userID uuid.UUID) ([]entity.Event, error) {
	if userID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	events, err := uc.repo.GetEventsForUser(ctx, userID)
	if err != nil {
		return nil, definitions.ErrEventNotFound
	}
	return events, nil
}

func (uc *UseCase) UpdateStatus(ctx context.Context, id uuid.UUID, status definitions.EventStatus) error {
	if id == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}
	return uc.repo.UpdateStatus(ctx, id, status)
}
