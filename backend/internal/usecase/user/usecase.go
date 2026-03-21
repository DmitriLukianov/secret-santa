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

func (uc *UseCase) Create(ctx context.Context, input dto.CreateUserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}

	user := entity.User{
		ID:    uuid.NewString(),
		Name:  input.Name,
		Email: input.Email,
	}

	return uc.repo.Create(ctx, user)
}

func (uc *UseCase) Get(ctx context.Context, id string) (*entity.User, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}

	return uc.repo.GetByID(ctx, id)
}
func (uc *UseCase) GetAll(ctx context.Context) ([]entity.User, error) {
	return uc.repo.GetAll(ctx)
}
func (uc *UseCase) Update(ctx context.Context, id string, input dto.UpdateUserInput) error {
	if input.Name == nil && input.Email == nil {
		return fmt.Errorf("nothing to update")
	}
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return uc.repo.Update(ctx, id, input.Name, input.Email)
}
func (uc *UseCase) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}

	return uc.repo.Delete(ctx, id)
}
