package auth

import (
	"context"
	"errors"
	"fmt"

	internalauth "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/usecase"
)

type UseCase struct {
	userUC usecase.UserUseCase
}

func New(userUC usecase.UserUseCase) *UseCase {
	return &UseCase{userUC: userUC}
}

// LoginWithOAuth — основная логика входа через OAuth
func (uc *UseCase) LoginWithOAuth(ctx context.Context, info internalauth.UserInfo) (string, error) {
	if info.ID == "" {
		return "", errors.New("oauth provider returned empty user id")
	}

	// 1. Сначала ищем существующего пользователя по oauth_provider + oauth_id
	user, err := uc.userUC.GetByOAuthID(ctx, info.ID, info.Provider)
	if err == nil && user != nil {
		// Пользователь уже есть — возвращаем его ID
		return user.ID.String(), nil
	}

	// 2. Если не нашли — создаём нового
	createInput := dto.CreateUserInput{
		Name:          info.Name,
		Email:         info.Email,
		OAuthID:       info.ID,
		OAuthProvider: info.Provider,
	}

	// Создаём пользователя и сразу получаем реальный сохранённый объект
	newUser, err := uc.userUC.Create(ctx, createInput)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// Возвращаем ID реально сохранённого пользователя
	return newUser.ID.String(), nil
}
