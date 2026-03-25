package assignment

import (
	"context"
	"fmt"
	"math/rand"

	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo            Repository
	participantRepo usecase.ParticipantRepository // теперь используем интерфейс из contracts.go
}

func New(repo Repository, participantRepo usecase.ParticipantRepository) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
	}
}

// Draw — запускает жеребьёвку (derangement: никто не дарит себе)
func (uc *UseCase) Draw(ctx context.Context, eventID uuid.UUID) error {
	if eventID == uuid.Nil {
		return fmt.Errorf("event id is required")
	}

	// Получаем участников
	participants, err := uc.participantRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}
	if len(participants) < 2 {
		return fmt.Errorf("not enough participants for drawing")
	}

	// Удаляем старую жеребьёвку
	if err := uc.repo.DeleteByEvent(ctx, eventID); err != nil {
		return fmt.Errorf("failed to delete old assignments: %w", err)
	}

	// Генерируем derangement
	assignments, err := uc.createDerangement(eventID, participants)
	if err != nil {
		return fmt.Errorf("failed to create derangement: %w", err)
	}

	// Сохраняем
	for _, a := range assignments {
		if err := uc.repo.Create(ctx, a); err != nil {
			return fmt.Errorf("failed to save assignment: %w", err)
		}
	}

	return nil
}

// createDerangement — алгоритм жеребьёвки без самоназначения
func (uc *UseCase) createDerangement(eventID uuid.UUID, participants []entity.Participant) ([]entity.Assignment, error) {
	n := len(participants)
	ids := make([]uuid.UUID, n)
	for i, p := range participants {
		ids[i] = p.UserID
	}

	for attempt := 0; attempt < 100; attempt++ {
		shuffled := make([]uuid.UUID, n)
		copy(shuffled, ids)
		rand.Shuffle(n, func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

		valid := true
		for i := 0; i < n; i++ {
			if shuffled[i] == ids[i] {
				valid = false
				break
			}
		}
		if valid {
			assignments := make([]entity.Assignment, n)
			for i := 0; i < n; i++ {
				assignments[i] = entity.NewAssignment(eventID, ids[i], shuffled[i])
			}
			return assignments, nil
		}
	}

	return nil, fmt.Errorf("failed to generate valid derangement after 100 attempts")
}

func (uc *UseCase) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error) {
	if eventID == uuid.Nil {
		return nil, fmt.Errorf("event id is required")
	}
	return uc.repo.GetByEvent(ctx, eventID)
}
