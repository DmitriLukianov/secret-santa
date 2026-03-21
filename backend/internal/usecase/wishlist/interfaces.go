package wishlist

import (
	"context"

	"secret-santa-backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, w entity.Wishlist) error

	GetByID(ctx context.Context, id string) (*entity.Wishlist, error)
	GetByUser(ctx context.Context, userID string) ([]entity.Wishlist, error)

	Update(
		ctx context.Context,
		id string,
		title, description, link, imageURL, visibility *string,
	) error

	Delete(ctx context.Context, id string) error
}
