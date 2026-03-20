package services

import (
	"context"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/repository"
)

type ParticipantService struct {
	repo repository.ParticipantRepository
}

func NewParticipantService(repo repository.ParticipantRepository) *ParticipantService {
	return &ParticipantService{repo: repo}
}

func (s *ParticipantService) JoinEvent(ctx context.Context, p domain.Participant) error {
	return s.repo.AddParticipant(ctx, p)
}

func (s *ParticipantService) GetParticipants(ctx context.Context, eventID string) ([]domain.Participant, error) {
	return s.repo.GetParticipantsByEvent(ctx, eventID)
}
func (s *ParticipantService) LeaveEvent(ctx context.Context, eventID, userID string) error {
	return s.repo.DeleteParticipant(ctx, eventID, userID)
}
