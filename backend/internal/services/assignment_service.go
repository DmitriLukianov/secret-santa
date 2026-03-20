package services

import (
	"context"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/repository"
)

type AssignmentService struct {
	repo repository.AssignmentRepository
}

func NewAssignmentService(repo repository.AssignmentRepository) *AssignmentService {
	return &AssignmentService{repo: repo}
}

func (s *AssignmentService) Create(ctx context.Context, a domain.Assignment) error {
	return s.repo.Create(ctx, a)
}
func (s *AssignmentService) GetMy(ctx context.Context, userID string) (*domain.Assignment, error) {
	return s.repo.GetByGiver(ctx, userID)
}
