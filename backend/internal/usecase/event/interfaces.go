package event

import (
	"context"
	"secret-santa-backend/internal/entity"
)

// Repository — интерфейс, который должен реализовать репозиторий
type Repository interface {
	Create(ctx context.Context, event entity.Event) (entity.Event, error)
	GetByID(ctx context.Context, id string) (entity.Event, error)
	List(ctx context.Context) ([]entity.Event, error)
	Update(ctx context.Context, event entity.Event) (entity.Event, error)
	Delete(ctx context.Context, id string) error
}
