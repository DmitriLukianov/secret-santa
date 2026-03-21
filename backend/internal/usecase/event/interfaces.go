package event

import (
	"context"
	"secret-santa-backend/internal/entity"
)

type Repository interface {
	Create(ctx context.Context, event entity.Event) error
	GetByID(ctx context.Context, id string) (*entity.Event, error)
	GetAll(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, id string, name, description *string) error
	Delete(ctx context.Context, id string) error
}
