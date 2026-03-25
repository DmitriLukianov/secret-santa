package auth

import (
	"context"
	"errors"
	"fmt"

	internalauth "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"
)

type UseCase struct {
	userUC usecase.UserUseCase // теперь используем публичный интерфейс
}

func New(userUC usecase.UserUseCase) *UseCase {
	return &UseCase{userUC: userUC}
}

func (uc *UseCase) LoginWithOAuth(ctx context.Context, info internalauth.UserInfo) (string, error) {
	if info.ID == "" {
		return "", errors.New("oauth provider returned empty user id")
	}

	// 1. Ищем существующего пользователя
	existing, err := uc.userUC.GetByOAuthID(ctx, info.ID, info.Provider)
	if err == nil && existing != nil {
		return existing.ID.String(), nil
	}

	// 2. Создаём нового
	newUser := entity.NewUser(info.Name, info.Email, info.ID, info.Provider)

	_, err = uc.userUC.Create(ctx, dto.CreateUserInput{ // используем DTO
		Name:          info.Name,
		Email:         info.Email,
		OAuthID:       info.ID,
		OAuthProvider: info.Provider,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return newUser.ID.String(), nil
}
