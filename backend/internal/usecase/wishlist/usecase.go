package wishlist

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (uc *UseCase) Create(ctx context.Context, participantID uuid.UUID, visibility string) (entity.Wishlist, error) {
	if uc.log != nil {
		uc.log.Info("create wishlist started",
			slog.String("participant_id", participantID.String()),
			slog.String("visibility", visibility),
		)
	}

	if participantID == uuid.Nil {
		return entity.Wishlist{}, definitions.ErrInvalidUserInput
	}

	wishlist := entity.NewWishlist(participantID, visibility)

	createdWishlist, err := uc.repo.Create(ctx, wishlist)
	if err != nil {
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
			slog.String("wishlist_id", createdWishlist.ID.String()),
			slog.String("participant_id", participantID.String()),
		)
	}

	return createdWishlist, nil
}

func (uc *UseCase) AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	if uc.log != nil {
		uc.log.Info("add wishlist item started",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("title", title),
		)
	}

	if wishlistID == uuid.Nil {
		return entity.WishlistItem{}, definitions.ErrInvalidUserInput
	}

	item := entity.NewWishlistItem(wishlistID, title, link, imageURL, comment)

	createdItem, err := uc.repo.CreateItem(ctx, item)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to add wishlist item",
				slog.String("wishlist_id", wishlistID.String()),
				slog.String("error", err.Error()),
			)
		}
		return entity.WishlistItem{}, fmt.Errorf("failed to add item: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("wishlist item added successfully",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("item_id", createdItem.ID.String()),
		)
	}

	return createdItem, nil
}

func (uc *UseCase) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	if participantID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	if uc.log != nil {
		uc.log.Info("get wishlist by participant started", slog.String("participant_id", participantID.String()))
	}
	return uc.repo.GetByParticipant(ctx, participantID)
}

func (uc *UseCase) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	if wishlistID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	if uc.log != nil {
		uc.log.Info("get wishlist items started", slog.String("wishlist_id", wishlistID.String()))
	}
	return uc.repo.GetItems(ctx, wishlistID)
}

func (uc *UseCase) GetForUser(ctx context.Context, eventID, participantID, requesterID uuid.UUID) (*entity.Wishlist, error) {
	if eventID == uuid.Nil || participantID == uuid.Nil || requesterID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
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
		if uc.log != nil {
			uc.log.Error("failed to get wishlist by participant",
				slog.String("participant_id", participantID.String()),
				slog.String("error", err.Error()),
			)
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", definitions.ErrWishlistNotFound, err)
		}
		return nil, fmt.Errorf("wishlist not found: %w", err)
	}

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
		return nil, definitions.ErrNotSanta
	}

	if uc.log != nil {
		uc.log.Info("wishlist access granted to santa",
			slog.String("requester_id", requesterID.String()),
			slog.String("receiver_user_id", participant.UserID.String()),
		)
	}

	return wishlist, nil
}

func (uc *UseCase) UpdateItem(ctx context.Context, itemID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	if uc.log != nil {
		uc.log.Info("update wishlist item started",
			slog.String("item_id", itemID.String()),
			slog.String("title", title),
		)
	}

	if itemID == uuid.Nil {
		return entity.WishlistItem{}, definitions.ErrInvalidUserInput
	}

	if err := uc.repo.UpdateItem(ctx, itemID, title, link, imageURL, comment); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update wishlist item", slog.String("error", err.Error()))
		}
		return entity.WishlistItem{}, fmt.Errorf("failed to update item: %w", err)
	}

	// Получаем актуальный item после обновления
	itemPtr, err := uc.repo.GetItemByID(ctx, itemID)
	if err != nil {
		return entity.WishlistItem{}, fmt.Errorf("failed to get updated item: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("wishlist item updated successfully", slog.String("item_id", itemID.String()))
	}

	return *itemPtr, nil // ← исправление здесь
}

func (uc *UseCase) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	if uc.log != nil {
		uc.log.Info("delete wishlist item started", slog.String("item_id", itemID.String()))
	}

	if itemID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}

	if err := uc.repo.DeleteItem(ctx, itemID); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to delete wishlist item", slog.String("error", err.Error()))
		}
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if uc.log != nil {
		uc.log.Info("wishlist item deleted successfully", slog.String("item_id", itemID.String()))
	}
	return nil
}

func (uc *UseCase) GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.WishlistItem, error) {
	if itemID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	if uc.log != nil {
		uc.log.Info("get wishlist item by id started", slog.String("item_id", itemID.String()))
	}

	return uc.repo.GetItemByID(ctx, itemID)
}

func (uc *UseCase) GetByID(ctx context.Context, wishlistID uuid.UUID) (*entity.Wishlist, error) {
	if wishlistID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	if uc.log != nil {
		uc.log.Info("get wishlist by id started", slog.String("wishlist_id", wishlistID.String()))
	}

	return uc.repo.GetByID(ctx, wishlistID)
}
