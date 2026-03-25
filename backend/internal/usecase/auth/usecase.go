package auth

import (
	"context"
	"errors"

	internalauth "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	userRepo UserRepository
}

func New(userRepo UserRepository) *UseCase {
	return &UseCase{userRepo: userRepo}
}

// LoginWithOAuth — найти пользователя по oauth_id или создать нового.
// Возвращает userID (UUID строкой).
func (uc *UseCase) LoginWithOAuth(ctx context.Context, info internalauth.UserInfo) (string, error) {
	if info.ID == "" {
		return "", errors.New("oauth provider returned empty user id")
	}

	// 1. Пробуем найти существующего пользователя по OAuthID
	existing, err := uc.userRepo.GetByOAuthID(ctx, info.ID)
	if err == nil && existing != nil {
		// пользователь уже есть — возвращаем его ID
		return existing.ID, nil
	}

	// 2. Если не нашли — создаём нового
	newUser := entity.User{
		ID:      uuid.NewString(),
		Name:    info.Name,
		Email:   info.Email,
		OAuthID: info.ID,
	}

	if err := uc.userRepo.Create(ctx, newUser); err != nil {
		return "", err
	}

	return newUser.ID, nil
}
