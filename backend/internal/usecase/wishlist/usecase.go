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
		)
	}
	return createdWishlist, nil
}

func (uc *UseCase) AddItem(ctx context.Context, wishlistID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	if wishlistID == uuid.Nil {
		return entity.WishlistItem{}, definitions.ErrInvalidUserInput
	}

	item := entity.NewWishlistItem(wishlistID, title, link, imageURL, comment)
	createdItem, err := uc.repo.CreateItem(ctx, item)
	if err != nil {
		if uc.log != nil {
			uc.log.Error("failed to add wishlist item", slog.String("error", err.Error()))
		}
		return entity.WishlistItem{}, fmt.Errorf("failed to add item: %w", err)
	}
	return createdItem, nil
}

func (uc *UseCase) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	if participantID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetByParticipant(ctx, participantID)
}

func (uc *UseCase) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	if wishlistID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetItems(ctx, wishlistID)
}

// GetForUser — получить вишлист с учётом visibility.
//
// Логика доступа:
//   - "public"     → виден всем авторизованным
//   - "friends"    → TODO: когда появится домен друзей; пока как "santa_only"
//   - "santa_only" → только Санта этого участника
//
// Владелец вишлиста всегда видит свой собственный.
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, definitions.ErrWishlistNotFound
		}
		return nil, fmt.Errorf("wishlist not found: %w", err)
	}

	participant, err := uc.participantRepo.GetByID(ctx, participantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get participant: %w", err)
	}

	// Владелец всегда видит свой вишлист
	if participant.UserID == requesterID {
		return wishlist, nil
	}

	switch wishlist.Visibility {
	case definitions.WishlistVisibilityPublic:
		// Виден всем
		return wishlist, nil

	case definitions.WishlistVisibilityFriends:
		// TODO: когда появится домен друзей — проверять дружбу.
		// Пока трактуем как santa_only (самый безопасный fallback).
		fallthrough

	case definitions.WishlistVisibilitySantaOnly:
		// Только Санта этого участника
		assignments, err := uc.assignmentRepo.GetByEvent(ctx, eventID)
		if err != nil {
			return nil, fmt.Errorf("failed to check assignment: %w", err)
		}

		for _, a := range assignments {
			if a.GiverID == requesterID && a.ReceiverID == participant.UserID {
				if uc.log != nil {
					uc.log.Info("wishlist access granted to santa",
						slog.String("requester_id", requesterID.String()),
					)
				}
				return wishlist, nil
			}
		}

		if uc.log != nil {
			uc.log.Warn("wishlist access denied",
				slog.String("requester_id", requesterID.String()),
				slog.String("visibility", wishlist.Visibility),
			)
		}
		return nil, definitions.ErrWishlistVisibilityForbidden

	default:
		return nil, definitions.ErrInvalidWishlistVisibility
	}
}

func (uc *UseCase) UpdateItem(ctx context.Context, itemID uuid.UUID, title string, link, imageURL, comment *string) (entity.WishlistItem, error) {
	if itemID == uuid.Nil {
		return entity.WishlistItem{}, definitions.ErrInvalidUserInput
	}

	if err := uc.repo.UpdateItem(ctx, itemID, title, link, imageURL, comment); err != nil {
		if uc.log != nil {
			uc.log.Error("failed to update wishlist item", slog.String("error", err.Error()))
		}
		return entity.WishlistItem{}, fmt.Errorf("failed to update item: %w", err)
	}

	itemPtr, err := uc.repo.GetItemByID(ctx, itemID)
	if err != nil {
		return entity.WishlistItem{}, fmt.Errorf("failed to get updated item: %w", err)
	}
	return *itemPtr, nil
}

func (uc *UseCase) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	if itemID == uuid.Nil {
		return definitions.ErrInvalidUserInput
	}
	if err := uc.repo.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (uc *UseCase) GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.WishlistItem, error) {
	if itemID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetItemByID(ctx, itemID)
}

func (uc *UseCase) GetByID(ctx context.Context, wishlistID uuid.UUID) (*entity.Wishlist, error) {
	if wishlistID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetByID(ctx, wishlistID)
}
