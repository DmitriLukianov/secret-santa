package assignment

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type UseCase struct {
	repo            Repository
	participantRepo ParticipantRepository
	eventRepo       EventRepository
	log             *slog.Logger
}

func New(repo Repository, participantRepo ParticipantRepository, eventRepo EventRepository) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
	}
}

func NewWithLogger(repo Repository, participantRepo ParticipantRepository, eventRepo EventRepository, log *slog.Logger) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
		log:             log,
	}
}

// Draw — запускает жеребьёвку (только организатор)
func (uc *UseCase) Draw(ctx context.Context, eventID, userID uuid.UUID) error {
	if uc.log != nil {
		uc.log.Info("draw started",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
		)
	}

	if eventID == uuid.Nil || userID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	// 1. Получаем событие
	eventPtr, err := uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get event", slog.String("error", err.Error()))
		}
		return fmt.Errorf("%w: %w", definitions.ErrEventNotFound, err)
	}

	// 2. Проверяем, что это организатор
	if eventPtr.OrganizerID != userID {
		return definitions.ErrNotOrganizer
	}

	// 3. Проверяем статус события
	if !eventPtr.IsDrawable() {
		if uc.log != nil {
			uc.log.Warn("draw not allowed due to status",
				slog.String("status", string(eventPtr.Status)),
			)
		}
		return definitions.ErrInvalidEventState
	}

	// 4. Получаем участников
	participants, err := uc.participantRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) < 3 {
		return definitions.ErrNotEnoughParticipants
	}

	// 5. Генерируем новую жеребьёвку (derangement)
	assignments, err := uc.createDerangement(eventID, participants)
	if err != nil {
		return fmt.Errorf("failed to create derangement: %w", err)
	}

	// FIXED: вся жеребьёвка теперь в одной атомарной транзакции
	if err := uc.repo.TransactionalDraw(ctx, eventID, assignments, entity.EventStatusDrawingDone); err != nil {
		return fmt.Errorf("failed to execute draw transaction: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("draw completed successfully",
			slog.String("event_id", eventID.String()),
			slog.Int("assignments_created", len(assignments)),
		)
	}

	return nil
}

// createDerangement — улучшенный алгоритм (никто не дарит себе)
func (uc *UseCase) createDerangement(eventID uuid.UUID, participants []entity.Participant) ([]entity.Assignment, error) {
	n := len(participants)
	ids := make([]uuid.UUID, n)
	for i, p := range participants {
		ids[i] = p.UserID
	}

	maxAttempts := 200 // FIXED: увеличено для большей надёжности

	for attempt := 0; attempt < maxAttempts; attempt++ {
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

	return nil, fmt.Errorf("failed to generate valid derangement after %d attempts", maxAttempts)
}

// GetByEvent — возвращает ТОЛЬКО свою пару (без изменений)
func (uc *UseCase) GetByEvent(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Assignment, error) {
	if eventID == uuid.Nil || userID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	assignments, err := uc.repo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}

	for _, a := range assignments {
		if a.GiverID == userID {
			return []entity.Assignment{a}, nil
		}
	}

	return []entity.Assignment{}, nil
}
