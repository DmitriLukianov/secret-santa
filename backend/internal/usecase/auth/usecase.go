package auth

import (
	"context"

	"secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *UseCase {
	return &UseCase{userRepo: userRepo}
}

func (uc *UseCase) LoginWithOAuth(ctx context.Context, user auth.UserInfo) (string, error) {
	// 🔍 пробуем найти пользователя
	existing, err := uc.userRepo.GetByID(ctx, user.ID)
	if err == nil && existing != nil {
		return existing.ID, nil
	}

	// ➕ если нет — создаём
	newUser := entity.User{
		ID:    uuid.NewString(),
		Name:  user.Name,
		Email: user.Email,
	}

	err = uc.userRepo.Create(ctx, newUser)
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
