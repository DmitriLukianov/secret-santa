package wishlist

import (
	"context"
	"errors"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Create(ctx context.Context, input dto.CreateWishlistInput) error {
	if input.UserID == "" {
		return errors.New("user_id is required")
	}

	w := entity.Wishlist{
		ID:          uuid.NewString(),
		UserID:      input.UserID,
		Title:       input.Title,
		Description: input.Description,
		Link:        input.Link,
		ImageURL:    input.ImageURL,
		Visibility:  input.Visibility,
	}

	return uc.repo.Create(ctx, w)
}

func (uc *UseCase) Get(ctx context.Context, id string) (*entity.Wishlist, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByUser(ctx context.Context, userID string) ([]entity.Wishlist, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	return uc.repo.GetByUser(ctx, userID)
}

func (uc *UseCase) Update(ctx context.Context, id string, input dto.UpdateWishlistInput) error {
	if id == "" {
		return errors.New("id is required")
	}

	return uc.repo.Update(
		ctx,
		id,
		input.Title,
		input.Description,
		input.Link,
		input.ImageURL,
		input.Visibility,
	)
}

func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return uc.repo.Delete(ctx, id)
}
