package invitation

import (
	"context"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo          Repository
	eventUC       usecase.EventUseCase
	participantUC usecase.ParticipantUseCase
	log           *slog.Logger
}

func New(repo Repository, eventUC usecase.EventUseCase, participantUC usecase.ParticipantUseCase) *UseCase {
	return &UseCase{repo: repo, eventUC: eventUC, participantUC: participantUC}
}

func NewWithLogger(repo Repository, eventUC usecase.EventUseCase, participantUC usecase.ParticipantUseCase, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, eventUC: eventUC, participantUC: participantUC, log: log}
}

func (uc *UseCase) GenerateInvite(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID) (dto.InvitationResponse, error) {
	if uc.log != nil {
		uc.log.Info("generate invitation started", slog.String("event_id", input.EventID.String()), slog.String("organizer_id", organizerID.String()))
	}

	event, err := uc.eventUC.GetByID(ctx, input.EventID)
	if err != nil {
		return dto.InvitationResponse{}, fmt.Errorf("%w: %w", definitions.ErrEventNotFound, err)
	}
	if event.OrganizerID != organizerID {
		return dto.InvitationResponse{}, definitions.ErrNotOrganizer
	}

	inv := entity.NewInvitation(input.EventID, organizerID, input.ExpiresIn)

	if err := uc.repo.Create(ctx, inv); err != nil {
		return dto.InvitationResponse{}, fmt.Errorf("failed to create invitation: %w", err)
	}

	inviteURL := fmt.Sprintf("https://yourdomain.com/invite/%s", inv.Token)

	if uc.log != nil {
		uc.log.Info("invitation generated", slog.String("token", inv.Token))
	}

	return dto.InvitationResponse{
		InviteURL: inviteURL,
		Token:     inv.Token,
		ExpiresAt: inv.ExpiresAt,
	}, nil
}

func (uc *UseCase) JoinByInvite(ctx context.Context, input dto.JoinByInvitationInput) error {
	if uc.log != nil {
		uc.log.Info("join by invitation started", slog.String("token", input.Token))
	}

	inv, err := uc.repo.GetByToken(ctx, input.Token)
	if err != nil {
		return definitions.ErrNotFound
	}

	// 1. Проверяем, что ссылка ещё действительна по времени
	if !inv.IsValid() {
		return definitions.ErrInvalidEventState
	}

	// 2. Проверяем статус события — можно присоединяться только пока открыт набор
	event, err := uc.eventUC.GetByID(ctx, inv.EventID)
	if err != nil {
		return err
	}

	if !event.CanAddParticipants() {
		return definitions.ErrInvalidEventState // ссылка закрывается после начала жеребьёвки
	}

	// 3. Добавляем участника
	_, err = uc.participantUC.Create(ctx, inv.EventID, input.UserID, entity.ParticipantRoleParticipant)
	if err != nil {
		return err
	}

	if uc.log != nil {
		uc.log.Info("user joined via invitation",
			slog.String("event_id", inv.EventID.String()),
			slog.String("user_id", input.UserID.String()))
	}

	return nil
}
