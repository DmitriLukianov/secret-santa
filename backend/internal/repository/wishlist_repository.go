package repository

import (
	"context"
	"secret-santa-backend/internal/domain"
)

type WishlistRepository interface {
	Create(ctx context.Context, wishlist domain.Wishlist) error
	GetByUser(ctx context.Context, eventID, userID string) (*domain.Wishlist, error)
	Update(ctx context.Context, id string, text *string) error
}
