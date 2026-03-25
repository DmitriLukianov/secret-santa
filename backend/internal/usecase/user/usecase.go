package user

import (
	"context"
	"fmt"

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

// Create
func (uc *UseCase) Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error) {
	if err := uc.validateCreateInput(input); err != nil {
		return entity.User{}, err
	}

	user := entity.NewUser(input.Name, input.Email, input.OAuthID, input.OAuthProvider)

	if err := uc.repo.Create(ctx, user); err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (uc *UseCase) validateCreateInput(input dto.CreateUserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if input.OAuthID == "" {
		return fmt.Errorf("oauthId is required")
	}
	if input.OAuthProvider == "" {
		return fmt.Errorf("oauthProvider is required")
	}
	return nil
}

// GetByID
func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is required")
	}
	return uc.repo.GetByID(ctx, id)
}

// GetByOAuthID
func (uc *UseCase) GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error) {
	if oauthID == "" || oauthProvider == "" {
		return nil, fmt.Errorf("oauthId and oauthProvider are required")
	}
	return uc.repo.GetByOAuthID(ctx, oauthID, oauthProvider)
}

// GetAll
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.User, error) {
	return uc.repo.GetAll(ctx)
}

// Update
func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	return uc.repo.Update(ctx, id, input.Name, input.Email)
}

// Delete
func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	return uc.repo.Delete(ctx, id)
}
