package wishlist

import (
	"context"

	"secret-santa-backend/internal/entity"
)

type Repository interface {
	// CREATE
	Create(ctx context.Context, w entity.Wishlist) error

	// READ
	GetByID(ctx context.Context, id string) (*entity.Wishlist, error)
	GetByUser(ctx context.Context, userID string) ([]entity.Wishlist, error)

	// UPDATE
	Update(
		ctx context.Context,
		id string,
		title, description, link, imageURL, visibility *string,
	) error

	// DELETE
	Delete(ctx context.Context, id string) error
}
