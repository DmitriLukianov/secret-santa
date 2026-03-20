package repository

import (
	"context"
	"secret-santa-backend/internal/domain"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event domain.Event) error
	GetEventByID(ctx context.Context, id string) (*domain.Event, error)
	GetEvents(ctx context.Context) ([]domain.Event, error)
	UpdateEvent(ctx context.Context, id string, name, description *string) error
	DeleteEvent(ctx context.Context, id string) error
}
