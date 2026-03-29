package auth

import (
	"context"
	"errors"
	"fmt"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/oauth"
	"secret-santa-backend/internal/usecase"

	"log/slog"
)

type UseCase struct {
	userUC usecase.UserUseCase
	log    *slog.Logger
}

func New(userUC usecase.UserUseCase) *UseCase {
	return &UseCase{userUC: userUC}
}

func NewWithLogger(userUC usecase.UserUseCase, log *slog.Logger) *UseCase {
	return &UseCase{userUC: userUC, log: log}
}

// LoginWithOAuth — основная логика OAuth (чисто по boilerplate)
func (uc *UseCase) LoginWithOAuth(ctx context.Context, info oauth.UserInfo) (string, error) {
	if uc.log != nil {
		uc.log.Info("oauth login started",
			slog.String("provider", info.Provider),
			slog.String("oauth_id", info.ID),
		)
	}

	if info.ID == "" {
		return "", definitions.ErrMissingOAuthCode
	}

	// 1. Ищем существующего пользователя
	user, err := uc.userUC.GetByOAuthID(ctx, info.ID, info.Provider)
	if err == nil && user != nil {
		if uc.log != nil {
			uc.log.Info("oauth user found", slog.String("user_id", user.ID.String()))
		}
		return user.ID.String(), nil
	}

	if err != nil && !errors.Is(err, definitions.ErrUserNotFound) {
		return "", fmt.Errorf("failed to lookup oauth user: %w", err)
	}

	// 2. Создаём нового пользователя
	createInput := dto.CreateUserInput{
		Name:          info.Name,
		Email:         info.Email,
		OAuthID:       info.ID,
		OAuthProvider: info.Provider,
	}

	if _, err = uc.userUC.Create(ctx, createInput); err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	// 3. Получаем реальный пользователь из БД (гарантируем правильный ID)
	savedUser, err := uc.userUC.GetByOAuthID(ctx, info.ID, info.Provider)
	if err != nil {
		return "", fmt.Errorf("failed to get saved user after creation: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("new oauth user created", slog.String("user_id", savedUser.ID.String()))
	}

	return savedUser.ID.String(), nil
}
