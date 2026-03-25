package wishlist

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

// Create — создаёт вишлист для участника
func (uc *UseCase) Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error) {
	wishlist := entity.NewWishlist(participantID, visibility)

	if err := uc.repo.Create(ctx, wishlist); err != nil {
		return entity.Wishlist{}, fmt.Errorf("failed to create wishlist: %w", err)
	}

	return wishlist, nil
}

// AddItem — добавляет элемент в вишлист
func (uc *UseCase) AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	item := entity.NewWishlistItem(wishlistID, title, link, imageURL, comment)

	if err := uc.repo.CreateItem(ctx, item); err != nil {
		return entity.WishlistItem{}, fmt.Errorf("failed to add item: %w", err)
	}

	return item, nil
}

func (uc *UseCase) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	if participantID == uuid.Nil {
		return nil, fmt.Errorf("participant id is required")
	}
	return uc.repo.GetByParticipant(ctx, participantID)
}

func (uc *UseCase) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	if wishlistID == uuid.Nil {
		return nil, fmt.Errorf("wishlist id is required")
	}
	return uc.repo.GetItems(ctx, wishlistID)
}
