package user

import (
	"context"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, log: log}
}

func (uc *UseCase) Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error) {
	if uc.log != nil {
		uc.log.Info("create user started",
			slog.String("email", input.Email),
			slog.String("oauth_provider", input.OAuthProvider),
		)
	}

	user := entity.NewUser(input.Name, input.Email, input.OAuthID, input.OAuthProvider)

	createdUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create user",
				slog.String("email", input.Email),
				slog.String("error", err.Error()),
			)
		}
		return entity.User{}, err
	}

	if uc.log != nil {
		uc.log.Info("user created successfully",
			slog.String("user_id", createdUser.ID.String()),
			slog.String("email", createdUser.Email),
		)
	}

	return createdUser, nil
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, definitions.ErrInvalidUUID
	}
	if uc.log != nil {
		uc.log.Info("get user by id started", slog.String("user_id", id.String()))
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error) {
	if oauthID == "" || oauthProvider == "" {
		return nil, definitions.ErrInvalidUserInput
	}
	if uc.log != nil {
		uc.log.Info("get user by oauth started",
			slog.String("oauth_id", oauthID),
			slog.String("oauth_provider", oauthProvider),
		)
	}
	return uc.repo.GetByOAuthID(ctx, oauthID, oauthProvider)
}

func (uc *UseCase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, definitions.ErrInvalidUserInput
	}
	if uc.log != nil {
		uc.log.Info("get user by email started", slog.String("email", email))
	}
	return uc.repo.GetByEmail(ctx, email)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error {
	if id == uuid.Nil {
		return definitions.ErrInvalidUUID
	}
	if uc.log != nil {
		uc.log.Info("update user started", slog.String("user_id", id.String()))
	}
	if err := uc.repo.Update(ctx, id, input.Name, input.Email); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update user", slog.String("user_id", id.String()), slog.String("error", err.Error()))
		}
		return err
	}
	if uc.log != nil {
		uc.log.Info("user updated successfully", slog.String("user_id", id.String()))
	}
	return nil
}

