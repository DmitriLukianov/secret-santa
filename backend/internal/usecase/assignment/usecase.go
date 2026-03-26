package assignment

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo            Repository
	participantRepo usecase.ParticipantRepository
	eventRepo       usecase.EventUseCase
	log             *slog.Logger
}

func New(repo Repository, participantRepo usecase.ParticipantRepository, eventRepo usecase.EventUseCase) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
	}
}

func NewWithLogger(repo Repository, participantRepo usecase.ParticipantRepository, eventRepo usecase.EventUseCase, log *slog.Logger) *UseCase {
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
		return fmt.Errorf("event id and user id are required")
	}

	// 1. Получаем событие
	event, err := uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get event",
				slog.String("event_id", eventID.String()),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to get event: %w", err)
	}

	// 2. Проверяем, что это организатор
	if event.OrganizerID != userID {
		if uc.log != nil {
			uc.log.Warn("unauthorized draw attempt",
				slog.String("event_id", eventID.String()),
				slog.String("user_id", userID.String()),
				slog.String("organizer_id", event.OrganizerID.String()),
			)
		}
		return fmt.Errorf("only the event organizer can start the draw")
	}

	// 3. Проверяем статус события
	if !event.IsDrawable() {
		if uc.log != nil {
			uc.log.Warn("draw not allowed due to status",
				slog.String("event_id", eventID.String()),
				slog.String("status", string(event.Status)),
			)
		}
		return fmt.Errorf("drawing already performed or event is not in draft status")
	}

	// 4. Получаем участников
	participants, err := uc.participantRepo.GetByEvent(ctx, eventID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get participants",
				slog.String("event_id", eventID.String()),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) < 3 {
		if uc.log != nil {
			uc.log.Warn("not enough participants",
				slog.String("event_id", eventID.String()),
				slog.Int("count", len(participants)),
			)
		}
		return fmt.Errorf("not enough participants for drawing (minimum 3 required)")
	}

	if uc.log != nil {
		uc.log.Info("starting derangement", slog.Int("participants_count", len(participants)))
	}

	// 5. Удаляем старую жеребьёвку
	if err := uc.repo.DeleteByEvent(ctx, eventID); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete old assignments",
				slog.String("event_id", eventID.String()),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to delete old assignments: %w", err)
	}

	// 6. Генерируем новую жеребьёвку
	assignments, err := uc.createDerangement(eventID, participants)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create derangement",
				slog.String("event_id", eventID.String()),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to create derangement: %w", err)
	}

	// 7. Сохраняем новые пары
	for _, a := range assignments {
		if err := uc.repo.Create(ctx, a); err != nil {
			if uc.log != nil {
				uc.log.Error("failed to save assignment",
					slog.String("event_id", eventID.String()),
					slog.String("giver_id", a.GiverID.String()),
					slog.String("receiver_id", a.ReceiverID.String()),
					slog.String("error", err.Error()),
				)
			}
			return fmt.Errorf("failed to save assignment: %w", err)
		}
	}

	// 8. Меняем статус события на "active"
	event.MarkAsDrawn()
	status := string(entity.EventStatusActive)
	if err := uc.eventRepo.Update(ctx, eventID, dto.UpdateEventInput{Status: &status}); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update event status",
				slog.String("event_id", eventID.String()),
				slog.String("status", status),
				slog.String("error", err.Error()),
			)
		}
		return fmt.Errorf("failed to update event status: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("draw completed successfully",
			slog.String("event_id", eventID.String()),
			slog.Int("assignments_created", len(assignments)),
		)
	}

	return nil
}

// createDerangement — алгоритм derangement (никто не дарит себе)
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

// GetByEvent — возвращает ТОЛЬКО свою пару (даже организатор видит только свою)
func (uc *UseCase) GetByEvent(ctx context.Context, eventID, userID uuid.UUID) ([]entity.Assignment, error) {
	if eventID == uuid.Nil || userID == uuid.Nil {
		return nil, fmt.Errorf("event id and user id are required")
	}

	assignments, err := uc.repo.GetByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignments: %w", err)
	}

	// Ищем только свою пару (где пользователь — giver)
	for _, a := range assignments {
		if a.GiverID == userID {
			if uc.log != nil {
				uc.log.Info("GetByEvent: returned own assignment",
					slog.String("event_id", eventID.String()),
					slog.String("user_id", userID.String()),
					slog.String("receiver_id", a.ReceiverID.String()),
				)
			}
			return []entity.Assignment{a}, nil
		}
	}

	// Если пользователь не является giver ни в одной паре
	if uc.log != nil {
		uc.log.Info("GetByEvent: no assignment found for user",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", userID.String()),
		)
	}
	return []entity.Assignment{}, nil
}
