package services

import (
	"context"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/repository"
)

type WishlistService struct {
	repo repository.WishlistRepository
}

func NewWishlistService(repo repository.WishlistRepository) *WishlistService {
	return &WishlistService{repo: repo}
}

func (s *WishlistService) Create(ctx context.Context, wishlist domain.Wishlist) error {
	return s.repo.Create(ctx, wishlist)
}

func (s *WishlistService) GetByUser(ctx context.Context, eventID, userID string) (*domain.Wishlist, error) {
	return s.repo.GetByUser(ctx, eventID, userID)
}

func (s *WishlistService) Update(ctx context.Context, id string, text *string) error {
	return s.repo.Update(ctx, id, text)
}
