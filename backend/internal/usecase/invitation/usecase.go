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
	repo           Repository
	eventRepo      EventRepository
	participantUC  usecase.ParticipantUseCase
	emailService   usecase.EmailService
	notificationUC usecase.NotificationUseCase
	baseURL        string
	log            *slog.Logger
}

func (uc *UseCase) SetNotificationUC(notificationUC usecase.NotificationUseCase) {
	uc.notificationUC = notificationUC
}

func New(repo Repository, eventRepo EventRepository, participantUC usecase.ParticipantUseCase, baseURL string) *UseCase {
	return &UseCase{
		repo:          repo,
		eventRepo:     eventRepo,
		participantUC: participantUC,
		baseURL:       baseURL,
	}
}

func NewWithLogger(repo Repository, eventRepo EventRepository, participantUC usecase.ParticipantUseCase, emailService usecase.EmailService, baseURL string, log *slog.Logger) *UseCase {
	return &UseCase{
		repo:          repo,
		eventRepo:     eventRepo,
		participantUC: participantUC,
		emailService:  emailService,
		baseURL:       baseURL,
		log:           log,
	}
}

func (uc *UseCase) GenerateInvite(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID) (dto.InvitationResponse, error) {
	if uc.log != nil {
		uc.log.Info("generate invitation started",
			slog.String("event_id", input.EventID.String()),
			slog.String("organizer_id", organizerID.String()),
		)
	}

	event, err := uc.eventRepo.GetByID(ctx, input.EventID)
	if err != nil {
		return dto.InvitationResponse{}, fmt.Errorf("%w: %s", definitions.ErrEventNotFound, err.Error())
	}

	if event.OrganizerID != organizerID {
		return dto.InvitationResponse{}, definitions.ErrNotOrganizer
	}

	inv := entity.NewInvitation(input.EventID, organizerID, input.ExpiresIn)
	createdInv, err := uc.repo.Create(ctx, inv)
	if err != nil {
		return dto.InvitationResponse{}, fmt.Errorf("failed to create invitation: %w", err)
	}

	inviteURL := fmt.Sprintf("%s/invite/%s", uc.baseURL, createdInv.Token)

	if uc.log != nil {
		uc.log.Info("invitation generated", slog.String("token", createdInv.Token))
	}

	return dto.InvitationResponse{
		InviteURL: inviteURL,
		Token:     createdInv.Token,
		ExpiresAt: createdInv.ExpiresAt,
	}, nil
}

func (uc *UseCase) SendEmailInvitation(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID, recipientEmail string) (dto.InvitationResponse, error) {
	resp, err := uc.GenerateInvite(ctx, input, organizerID)
	if err != nil {
		return dto.InvitationResponse{}, err
	}

	if uc.emailService != nil {
		event, err := uc.eventRepo.GetByID(ctx, input.EventID)
		if err != nil {
			return dto.InvitationResponse{}, fmt.Errorf("%w: %s", definitions.ErrEventNotFound, err.Error())
		}
		if err := uc.emailService.SendInvitationEmail(ctx, recipientEmail, event.Title, resp.InviteURL); err != nil {
			if uc.log != nil {
				uc.log.Warn("failed to send invitation email",
					slog.String("email", recipientEmail),
					slog.String("error", err.Error()),
				)
			}
		}
	}

	return resp, nil
}

func (uc *UseCase) JoinByInvite(ctx context.Context, input dto.JoinByInvitationInput) error {
	if uc.log != nil {
		uc.log.Info("join by invitation started", slog.String("token", input.Token))
	}

	inv, err := uc.repo.GetByToken(ctx, input.Token)
	if err != nil {
		return definitions.ErrNotFound
	}

	if !inv.IsValid() {
		return definitions.ErrInvalidEventState
	}

	event, err := uc.eventRepo.GetByID(ctx, inv.EventID)
	if err != nil {
		return err
	}

	if !event.CanAddParticipants() {
		return definitions.ErrInvalidEventState
	}

	_, err = uc.participantUC.Create(ctx, inv.EventID, input.UserID, definitions.ParticipantRoleParticipant)
	if err != nil {
		return err
	}

	if uc.notificationUC != nil {
		_ = uc.notificationUC.Notify(ctx, event.OrganizerID, "invitation_joined", map[string]string{
			"event_id": inv.EventID.String(),
			"user_id":  input.UserID.String(),
		})
	}

	if uc.log != nil {
		uc.log.Info("user joined via invitation",
			slog.String("event_id", inv.EventID.String()),
			slog.String("user_id", input.UserID.String()),
		)
	}
	return nil
}
