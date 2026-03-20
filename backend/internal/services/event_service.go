package services

import (
	"context"
	"fmt"
	"secret-santa-backend/internal/domain"
	"secret-santa-backend/internal/repository"
)

type EventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

func (s *EventService) CreateEvent(ctx context.Context, event domain.Event) error {

	if event.Name == "" {
		return fmt.Errorf("event name is required")
	}

	if event.OrganizerID == "" {
		return fmt.Errorf("organizer_id is required")
	}

	return s.eventRepo.CreateEvent(ctx, event)
}

func (s *EventService) GetEvent(ctx context.Context, id string) (*domain.Event, error) {
	return s.eventRepo.GetEventByID(ctx, id)
}

func (s *EventService) GetEvents(ctx context.Context) ([]domain.Event, error) {
	return s.eventRepo.GetEvents(ctx)
}

func (s *EventService) UpdateEvent(ctx context.Context, id string, name, description *string) error {
	return s.eventRepo.UpdateEvent(ctx, id, name, description)
}

func (s *EventService) DeleteEvent(ctx context.Context, id string) error {
	return s.eventRepo.DeleteEvent(ctx, id)
}
