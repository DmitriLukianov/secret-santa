package usecase

import (
	"context"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

// UserUseCase — публичный интерфейс
type UserUseCase interface {
	Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// EventUseCase — полный интерфейс событий (обновлённый)
type EventUseCase interface {
	Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status entity.EventStatus) error
	Delete(ctx context.Context, id uuid.UUID) error

	// === Работа со статусами ===
	OpenInvitation(ctx context.Context, id, userID uuid.UUID) error
	CloseRegistration(ctx context.Context, id, userID uuid.UUID) error
	StartDrawing(ctx context.Context, id, userID uuid.UUID) error
	Finish(ctx context.Context, id, userID uuid.UUID) error
	Cancel(ctx context.Context, id, userID uuid.UUID) error

	GetMyEvents(ctx context.Context, userID uuid.UUID) ([]entity.Event, error)
}

// ParticipantUseCase
type ParticipantUseCase interface {
	Create(ctx context.Context, eventID, userID uuid.UUID, role string) (entity.Participant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
	MarkGiftSent(ctx context.Context, participantID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error)
}

// WishlistUseCase
type WishlistUseCase interface {
	Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error)
	AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error)
	GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error)
	GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error)
	GetForUser(ctx context.Context, eventID, participantID, requesterID uuid.UUID) (*entity.Wishlist, error)
}

// AssignmentUseCase
type AssignmentUseCase interface {
	Draw(ctx context.Context, eventID, userID uuid.UUID) error
	GetByEvent(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Assignment, error)
}

// ParticipantRepository (для AssignmentUseCase)
type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}

// InvitationUseCase — публичный интерфейс приглашений
type InvitationUseCase interface {
	GenerateInvite(ctx context.Context, input dto.CreateInvitationInput, organizerID uuid.UUID) (dto.InvitationResponse, error)
	JoinByInvite(ctx context.Context, input dto.JoinByInvitationInput) error
}
