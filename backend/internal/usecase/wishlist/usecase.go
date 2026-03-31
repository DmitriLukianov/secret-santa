package wishlist

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/definitions" // ← добавлен
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5" // ← добавлен для ErrNoRows
)

type UseCase struct {
	repo            Repository
	participantRepo ParticipantRepository
	assignmentRepo  AssignmentRepository
	log             *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, participantRepo ParticipantRepository, assignmentRepo AssignmentRepository, log *slog.Logger) *UseCase {
	return &UseCase{
		repo:            repo,
		participantRepo: participantRepo,
		assignmentRepo:  assignmentRepo,
		log:             log,
	}
}

// Create — создаёт вишлист для участника
func (uc *UseCase) Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error) {
	if uc.log != nil {
		uc.log.Info("create wishlist started",
			slog.String("participant_id", participantID.String()),
			slog.String("visibility", visibility),
		)
	}

	wishlist := entity.NewWishlist(participantID, visibility)

	if err := uc.repo.Create(ctx, wishlist); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to create wishlist",
				slog.String("participant_id", participantID.String()),
				slog.String("error", err.Error()),
			)
		}
		return entity.Wishlist{}, fmt.Errorf("failed to create wishlist: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("wishlist created successfully",
			slog.String("wishlist_id", wishlist.ID.String()),
			slog.String("participant_id", participantID.String()),
		)
	}

	return wishlist, nil
}

// AddItem — добавляет элемент в вишлист
func (uc *UseCase) AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	if uc.log != nil {
		uc.log.Info("add wishlist item started",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("title", title),
			slog.Any("link", link),
			slog.Any("image_url", imageURL),
			slog.Any("comment", comment),
		)
	}

	item := entity.NewWishlistItem(wishlistID, title, link, imageURL, comment)

	if err := uc.repo.CreateItem(ctx, item); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to add wishlist item",
				slog.String("wishlist_id", wishlistID.String()),
				slog.String("title", title),
				slog.String("error", err.Error()),
			)
		}
		return entity.WishlistItem{}, fmt.Errorf("failed to add item: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("wishlist item added successfully",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("item_id", item.ID.String()),
		)
	}

	return item, nil
}

func (uc *UseCase) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	if participantID == uuid.Nil {
		return nil, fmt.Errorf("participant id is required")
	}
	if uc.log != nil {
		uc.log.Info("get wishlist by participant started", slog.String("participant_id", participantID.String()))
	}
	return uc.repo.GetByParticipant(ctx, participantID)
}

func (uc *UseCase) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	if wishlistID == uuid.Nil {
		return nil, fmt.Errorf("wishlist id is required")
	}
	if uc.log != nil {
		uc.log.Info("get wishlist items started", slog.String("wishlist_id", wishlistID.String()))
	}
	return uc.repo.GetItems(ctx, wishlistID)
}

// GetForUser — возвращает вишлист только если requester является Сантой
func (uc *UseCase) GetForUser(ctx context.Context, eventID, participantID, requesterID uuid.UUID) (*entity.Wishlist, error) {
	if eventID == uuid.Nil || participantID == uuid.Nil || requesterID == uuid.Nil {
		return nil, fmt.Errorf("eventID, participantID and requesterID are required")
	}

	if uc.log != nil {
		uc.log.Info("get wishlist for user started",
			slog.String("event_id", eventID.String()),
			slog.String("participant_id", participantID.String()),
			slog.String("requester_id", requesterID.String()),
		)
	}

	// 1. Получаем вишлист
	wishlist, err := uc.repo.GetByParticipant(ctx, participantID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get wishlist by participant",
				slog.String("participant_id", participantID.String()),
				slog.String("error", err.Error()),
			)
		}
		// 🔥 КРИТИЧЕСКИЙ ФИКС: теперь возвращаем 404 вместо 500
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", definitions.ErrWishlistNotFound, err)
		}
		return nil, fmt.Errorf("wishlist not found: %w", err)
	}

	// 2. Получаем участника (чтобы узнать реальный UserID получателя)
	participant, err := uc.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get participant by id",
				slog.String("participant_id", participantID.String()),
				slog.String("error", err.Error()),
			)
		}
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("participant found",
			slog.String("participant_id", participant.ID.String()),
			slog.String("user_id", participant.UserID.String()),
		)
	}

	// 3. Получаем все назначения события
	assignments, err := uc.assignmentRepo.GetByEvent(ctx, eventID)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to get assignments",
				slog.String("event_id", eventID.String()),
				slog.String("error", err.Error()),
			)
		}
		return nil, fmt.Errorf("failed to check assignment: %w", err)
	}

	// 4. Проверяем, является ли requester Сантой для этого получателя
	isSanta := false
	for _, a := range assignments {
		if a.GiverID == requesterID && a.ReceiverID == participant.UserID {
			isSanta = true
			break
		}
	}

	if !isSanta {
		if uc.log != nil {
			uc.log.Warn("wishlist access denied: not the santa",
				slog.String("requester_id", requesterID.String()),
				slog.String("receiver_user_id", participant.UserID.String()),
			)
		}
		return nil, fmt.Errorf("you are not the santa for this participant")
	}

	if uc.log != nil {
		uc.log.Info("wishlist access granted to santa",
			slog.String("requester_id", requesterID.String()),
			slog.String("receiver_user_id", participant.UserID.String()),
		)
	}

	return wishlist, nil
}
