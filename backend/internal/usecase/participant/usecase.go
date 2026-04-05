package participant

import (
	"context"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/helpers"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo           Repository
	notificationUC usecase.NotificationUseCase
	log            *slog.Logger
}

func (uc *UseCase) SetNotificationUC(notificationUC usecase.NotificationUseCase) {
	uc.notificationUC = notificationUC
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, log: log}
}

func (uc *UseCase) Create(ctx context.Context, eventID, userID uuid.UUID, role string) (entity.Participant, error) {
	if uc.log != nil {
		uc.log.Info("create participant started",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
			slog.String("role", role),
		)
	}

	if eventID == uuid.Nil || userID == uuid.Nil {
		return entity.Participant{}, definitions.ErrInvalidUserInput
	}

	participant := entity.NewParticipant(eventID, userID, role)
	created, err := uc.repo.Create(ctx, participant)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create participant", slog.String("error", err.Error()))
		}
		if helpers.IsDuplicateError(err) {
			return entity.Participant{}, definitions.ErrAlreadyParticipating
		}
		return entity.Participant{}, err
	}

	if uc.log != nil {
		uc.log.Info("participant created successfully",
			slog.String("participant_id", created.ID.String()),
		)
	}
	return created, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	if id == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error) {
	if eventID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetByEvent(ctx, eventID)
}

func (uc *UseCase) MarkGiftSent(ctx context.Context, participantID, requesterID uuid.UUID) error {
	if participantID == uuid.Nil || requesterID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	if uc.log != nil {
		uc.log.Info("mark gift sent started",
			slog.String("participant_id", participantID.String()),
			slog.String("requester_id", requesterID.String()),
		)
	}

	p, err := uc.repo.GetByID(ctx, participantID)
	if err != nil {
		return definitions.ErrParticipantNotFound
	}
	if p.UserID != requesterID {
		return definitions.ErrForbidden
	}

	if err := uc.repo.UpdateGiftSent(ctx, participantID, true); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to mark gift sent", slog.String("error", err.Error()))
		}
		return err
	}

	if uc.notificationUC != nil {
		_ = uc.notificationUC.Notify(ctx, p.UserID, "gift_sent", map[string]string{
			"participant_id": participantID.String(),
		})
	}

	if uc.log != nil {
		uc.log.Info("gift marked as sent", slog.String("participant_id", participantID.String()))
	}
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id, requesterID uuid.UUID) error {
	if id == uuid.Nil || requesterID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	if uc.log != nil {
		uc.log.Info("delete participant started",
			slog.String("participant_id", id.String()),
			slog.String("requester_id", requesterID.String()),
		)
	}

	p, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return definitions.ErrParticipantNotFound
	}

	if p.UserID != requesterID {
		return definitions.ErrForbidden
	}

	if err := uc.repo.Delete(ctx, id); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete participant", slog.String("error", err.Error()))
		}
		return err
	}

	if uc.log != nil {
		uc.log.Info("participant deleted successfully", slog.String("participant_id", id.String()))
	}
	return nil
}

func (uc *UseCase) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error) {
	if userID == uuid.Nil || eventID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetByUserAndEvent(ctx, userID, eventID)
}
