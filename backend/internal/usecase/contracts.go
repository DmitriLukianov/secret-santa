package usecase

import (
	"context"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

// UserUseCase — публичный интерфейс (используется в контроллерах и auth)
type UserUseCase interface {
	Create(ctx context.Context, input dto.CreateUserInput) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateUserInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// EventUseCase — публичный интерфейс для событий
type EventUseCase interface {
	Create(ctx context.Context, input dto.CreateEventInput, organizerID uuid.UUID) (entity.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error
	Delete(ctx context.Context, id uuid.UUID) error
}
