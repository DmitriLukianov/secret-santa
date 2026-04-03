package user

import (
	"context"
	"fmt"
	"log/slog"

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
			slog.Any("email", input.Email),
			slog.Any("oauth_provider", input.OAuthProvider),
		)
	}

	if err := uc.validateCreateInput(input); err != nil {
		if uc.log != nil {
			uc.log.Warn("create user validation failed",
				slog.Any("email", input.Email),
				slog.String("error", err.Error()),
			)
		}
		return entity.User{}, err
	}

	user := entity.NewUser(input.Name, input.Email, input.OAuthID, input.OAuthProvider)

	if err := uc.repo.Create(ctx, user); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create user",
				slog.String("email", user.Email),
				slog.String("error", err.Error()),
			)
		}
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("user created successfully",
			slog.String("user_id", user.ID.String()),
			slog.String("email", user.Email),
		)
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

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("id is required")
	}
	if uc.log != nil {
		uc.log.Info("get user by id started", slog.String("user_id", id.String()))
	}
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error) {
	if oauthID == "" || oauthProvider == "" {
		return nil, fmt.Errorf("oauthId and oauthProvider are required")
	}
	if uc.log != nil {
		uc.log.Info("get user by oauth started",
			slog.String("oauth_id", oauthID),
			slog.String("oauth_provider", oauthProvider),
		)
	}
	return uc.repo.GetByOAuthID(ctx, oauthID, oauthProvider)
}

func (uc *UseCase) GetAll(ctx context.Context) ([]entity.User, error) {
	if uc.log != nil {
		uc.log.Info("get all users started")
	}
	return uc.repo.GetAll(ctx)
}

func (uc *UseCase) Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	if uc.log != nil {
		uc.log.Info("update user started",
			slog.String("user_id", id.String()),
			slog.Any("name", input.Name),
			slog.Any("email", input.Email),
		)
	}
	if err := uc.repo.Update(ctx, id, input.Name, input.Email); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update user",
				slog.String("user_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	if uc.log != nil {
		uc.log.Info("user updated successfully", slog.String("user_id", id.String()))
	}
	return nil
}

func (uc *UseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("id is required")
	}
	if uc.log != nil {
		uc.log.Info("delete user started", slog.String("user_id", id.String()))
	}
	if err := uc.repo.Delete(ctx, id); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete user",
				slog.String("user_id", id.String()),
				slog.String("error", err.Error()),
			)
		}
		return err
	}
	if uc.log != nil {
		uc.log.Info("user deleted successfully", slog.String("user_id", id.String()))
	}
	return nil
}
