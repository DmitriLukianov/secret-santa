package usecase

import (
	"context"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

// UserUseCase — публичный интерфейс (используется в контроллерах и auth)
type UserUseCase interface {
	Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// EventUseCase — публичный интерфейс для событий
type EventUseCase interface {
	Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error
	Delete(ctx context.Context, id uuid.UUID) error
	Finish(ctx context.Context, id, userID uuid.UUID) error
}

// ParticipantUseCase — публичный интерфейс участников
type ParticipantUseCase interface {
	Create(ctx context.Context, eventID, userID uuid.UUID, role string) (entity.Participant, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
	MarkGiftSent(ctx context.Context, participantID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error)
}

// WishlistUseCase — публичный интерфейс вишлистов
type WishlistUseCase interface {
	Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error)
	AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error)
	GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error)
	GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error)
	GetForUser(ctx context.Context, eventID, participantID, requesterID uuid.UUID) (*entity.Wishlist, error)
}

// AssignmentUseCase — публичный интерфейс жеребьёвки
type AssignmentUseCase interface {
	Draw(ctx context.Context, eventID, userID uuid.UUID) error
	GetByEvent(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Assignment, error) // ← исправлено
}

// ParticipantRepository — минимальный интерфейс, который нужен AssignmentUseCase
type ParticipantRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error)
}
