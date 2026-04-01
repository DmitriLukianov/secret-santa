package wishlist

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, wishlist entity.Wishlist) error
	CreateItem(ctx context.Context, item entity.WishlistItem) error
	GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error)
	GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error)
	UpdateItem(ctx context.Context, itemID uuid.UUID, title string, link, imageURL, comment *string) error
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
}

// FIXED: добавили ParticipantRepository, чтобы получить UserID участника
// (нужно для корректной проверки "являешься ли ты Сантой")
type ParticipantRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error)
}

type AssignmentRepository interface {
	GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error)
}
