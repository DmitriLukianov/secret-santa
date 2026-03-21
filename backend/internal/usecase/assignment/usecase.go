package assignment

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	assignRepo AssignmentRepository
	partRepo   ParticipantRepository
}

func New(assignRepo AssignmentRepository, partRepo ParticipantRepository) *UseCase {
	return &UseCase{
		assignRepo: assignRepo,
		partRepo:   partRepo,
	}
}

// 🎁 GENERATE (жеребьёвка)
func (uc *UseCase) Draw(ctx context.Context, input dto.GenerateAssignmentInput) error {
	if input.EventID == "" {
		return errors.New("event_id is required")
	}

	// 1. Получаем участников
	participants, err := uc.partRepo.GetByEvent(ctx, input.EventID)
	if err != nil {
		return err
	}

	// 2. Проверка
	if len(participants) < 2 {
		return errors.New("not enough participants")
	}

	// 3. Проверка: уже есть жеребьёвка
	existing, err := uc.assignRepo.GetByEvent(ctx, input.EventID)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return errors.New("assignments already exist")
	}

	// 4. Перемешивание участников
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	r.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	// 5. Генерация пар (кольцевая схема)
	assignments := make([]entity.Assignment, 0, len(participants))

	for i := 0; i < len(participants); i++ {
		giver := participants[i]
		receiver := participants[(i+1)%len(participants)]

		// защита от "сам себе"
		if giver.UserID == receiver.UserID {
			return errors.New("invalid assignment: self assignment")
		}

		assignments = append(assignments, entity.Assignment{
			ID:         uuid.NewString(),
			EventID:    input.EventID,
			GiverID:    giver.UserID,
			ReceiverID: receiver.UserID,
		})
	}

	// 6. Сохраняем
	return uc.assignRepo.CreateMany(ctx, assignments)
}

// 📥 GET BY EVENT
func (uc *UseCase) GetByEvent(ctx context.Context, eventID string) ([]entity.Assignment, error) {
	if eventID == "" {
		return nil, errors.New("event_id is required")
	}

	return uc.assignRepo.GetByEvent(ctx, eventID)
}
