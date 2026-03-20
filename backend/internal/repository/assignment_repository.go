package repository

import (
	"context"
	"secret-santa-backend/internal/domain"
)

type AssignmentRepository interface {
	Create(ctx context.Context, a domain.Assignment) error
	GetByGiver(ctx context.Context, giverID string) (*domain.Assignment, error)
	GetByEvent(ctx context.Context, eventID string) ([]domain.Assignment, error)
}
