package chat

import (
	"context"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo            usecase.ChatRepository
	participantRepo usecase.ParticipantRepository
	assignmentRepo  usecase.AssignmentRepository
	log             *slog.Logger
}

func New(repo usecase.ChatRepository, participantRepo usecase.ParticipantRepository, assignmentRepo usecase.AssignmentRepository) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		assignmentRepo:  assignmentRepo,
	}
}

func NewWithLogger(repo usecase.ChatRepository, participantRepo usecase.ParticipantRepository, assignmentRepo usecase.AssignmentRepository, log *slog.Logger) *UseCase {
	uc := New(repo, participantRepo, assignmentRepo)
	uc.log = log
	return uc
}

// SendMessage — отправить сообщение (автоматически определяет пару)
func (uc *UseCase) SendMessage(ctx context.Context, eventID, userID uuid.UUID, content string) (entity.Message, error) {
	if content == "" {
		return entity.Message{}, definitions.ErrInvalidUserInput
	}

	if uc.log != nil {
		uc.log.Info("send message started",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
		)
	}

	assignments, err := uc.assignmentRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to get assignments: %w", err)
	}

	var receiverID uuid.UUID
	for _, a := range assignments {
		if a.GiverID == userID {
			receiverID = a.ReceiverID
			break
		}
		if a.ReceiverID == userID {
			receiverID = a.GiverID
			break
		}
	}
	if receiverID == uuid.Nil {
		return entity.Message{}, definitions.ErrNotSanta
	}

	msg := entity.NewMessage(eventID, userID, receiverID, content)

	// Теперь CreateMessage возвращает полную сущность из БД
	createdMsg, err := uc.repo.CreateMessage(ctx, msg)
	if err != nil {
		return entity.Message{}, fmt.Errorf("failed to create message: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("message sent successfully",
			slog.String("message_id", createdMsg.ID.String()),
		)
	}

	return createdMsg, nil
}

// GetRecipientChat — чат «Кому я Санта» (я — giver)
func (uc *UseCase) GetRecipientChat(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Message, error) {
	if uc.log != nil {
		uc.log.Info("get recipient chat started",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
		)
	}

	assignments, err := uc.assignmentRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}

	var receiverID uuid.UUID
	for _, a := range assignments {
		if a.GiverID == userID {
			receiverID = a.ReceiverID
			break
		}
	}
	if receiverID == uuid.Nil {
		return nil, definitions.ErrNotSanta
	}

	return uc.repo.GetMessagesByPair(ctx, eventID, userID, receiverID)
}

// GetSenderChat — чат «Кто мой Санта» (я — receiver)
func (uc *UseCase) GetSenderChat(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Message, error) {
	if uc.log != nil {
		uc.log.Info("get sender chat started",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
		)
	}

	assignments, err := uc.assignmentRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}

	var senderID uuid.UUID
	for _, a := range assignments {
		if a.ReceiverID == userID {
			senderID = a.GiverID
			break
		}
	}
	if senderID == uuid.Nil {
		return nil, definitions.ErrNotSanta
	}

	return uc.repo.GetMessagesByPair(ctx, eventID, senderID, userID)
}
