package assignment

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo            Repository
	participantRepo ParticipantRepository
	eventRepo       EventRepository
	userUC          usecase.UserUseCase
	emailService    usecase.EmailService
	log             *slog.Logger
}

func New(repo Repository, participantRepo ParticipantRepository, eventRepo EventRepository, userUC usecase.UserUseCase) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
		userUC:          userUC,
	}
}

func NewWithLogger(
	repo Repository,
	participantRepo ParticipantRepository,
	eventRepo EventRepository,
	userUC usecase.UserUseCase,
	emailService usecase.EmailService,
	log *slog.Logger,
) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		eventRepo:       eventRepo,
		userUC:          userUC,
		emailService:    emailService,
		log:             log,
	}
}

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

	eventPtr, err := uc.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get event", slog.String("error", err.Error()))
		}
		return fmt.Errorf("%w: %s", definitions.ErrEventNotFound, err.Error())
	}

	if eventPtr.OrganizerID != userID {
		return definitions.ErrNotOrganizer
	}

	if !eventPtr.IsDrawable() {
		if uc.log != nil {
			uc.log.Warn("draw not allowed due to status",
				slog.String("status", string(eventPtr.Status)),
			)
		}
		return definitions.ErrInvalidEventState
	}

	participants, err := uc.participantRepo.GetByEvent(ctx, eventID)
	if err != nil {
		return fmt.Errorf("failed to get participants: %w", err)
	}

	if len(participants) < 3 {
		return definitions.ErrNotEnoughParticipants
	}

	assignments, err := uc.createDerangement(eventID, participants)
	if err != nil {
		return fmt.Errorf("failed to create derangement: %w", err)
	}

	if err := uc.repo.TransactionalDraw(ctx, eventID, assignments, definitions.EventStatusDrawingDone); err != nil {
		return fmt.Errorf("failed to execute draw transaction: %w", err)
	}

	if uc.emailService != nil && uc.userUC != nil {
		notified := 0
		for _, p := range participants {
			userPtr, err := uc.userUC.GetByID(ctx, p.UserID)
			if err != nil || userPtr == nil {
				if uc.log != nil {
					uc.log.Warn("failed to resolve participant for draw notification",
						slog.String("user_id", p.UserID.String()),
						slog.String("error", fmt.Sprint(err)),
					)
				}
				continue
			}
			if err := uc.emailService.SendDrawNotification(ctx, userPtr.Email, eventPtr.Title); err != nil {
				if uc.log != nil {
					uc.log.Warn("failed to send draw notification",
						slog.String("user_id", userPtr.ID.String()),
						slog.String("email", userPtr.Email),
						slog.String("error", err.Error()),
					)
				}
				continue
			}
			notified++
		}
		if uc.log != nil {
			uc.log.Info("draw notifications sent", slog.Int("notified_users", notified))
		}
	}

	if uc.log != nil {
		uc.log.Info("draw completed successfully",
			slog.String("event_id", eventID.String()),
			slog.Int("assignments_created", len(assignments)),
		)
	}
	return nil
}

func (uc *UseCase) createDerangement(eventID uuid.UUID, participants []entity.Participant) ([]entity.Assignment, error) {
	n := len(participants)
	ids := make([]uuid.UUID, n)
	for i, p := range participants {
		ids[i] = p.UserID
	}

	const maxAttempts = 200
	for attempt := 0; attempt < maxAttempts; attempt++ {
		shuffled := make([]uuid.UUID, n)
		copy(shuffled, ids)

		if err := cryptoShuffle(shuffled); err != nil {
			return nil, fmt.Errorf("failed to shuffle: %w", err)
		}

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

func cryptoShuffle(s []uuid.UUID) error {
	n := len(s)
	for i := n - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(jBig.Int64())
		s[i], s[j] = s[j], s[i]
	}
	return nil
}

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
