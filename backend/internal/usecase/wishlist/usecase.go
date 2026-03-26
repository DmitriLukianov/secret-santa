package wishlist

import (
	"context"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/entity"
	"secret-santa-backend/internal/usecase"

	"github.com/google/uuid"
)

type UseCase struct {
	repo         Repository
	assignmentUC usecase.AssignmentUseCase // ← добавили
	log          *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, assignmentUC usecase.AssignmentUseCase, log *slog.Logger) *UseCase {
	return &UseCase{
		repo:         repo,
		assignmentUC: assignmentUC,
		log:          log,
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

// GetForUser — возвращает вишлист только если requester является Сантой этого участника
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

	wishlist, err := uc.repo.GetByParticipant(ctx, participantID)
	if err != nil {
		return nil, fmt.Errorf("wishlist not found: %w", err)
	}

	// Проверяем, является ли requester Сантой этого участника
	assignments, err := uc.assignmentUC.GetByEvent(ctx, eventID, requesterID)
	if err != nil {
		return nil, fmt.Errorf("failed to check assignment: %w", err)
	}

	isSanta := false
	for _, a := range assignments {
		if a.ReceiverID == wishlist.ParticipantID {
			isSanta = true
			break
		}
	}

	if !isSanta {
		if uc.log != nil {
			uc.log.Warn("wishlist access denied: not the santa",
				slog.String("requester_id", requesterID.String()),
				slog.String("participant_id", participantID.String()),
			)
		}
		return nil, fmt.Errorf("you are not the santa for this participant")
	}

	if uc.log != nil {
		uc.log.Info("wishlist access granted",
			slog.String("requester_id", requesterID.String()),
			slog.String("participant_id", participantID.String()),
		)
	}

	return wishlist, nil
}
