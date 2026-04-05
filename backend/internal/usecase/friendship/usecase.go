package friendship

import (
	"context"
	"errors"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/google/uuid"
)

type UseCase struct {
	repo Repository
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) SendRequest(ctx context.Context, requesterID, addresseeID uuid.UUID) (entity.Friendship, error) {
	if requesterID == uuid.Nil || addresseeID == uuid.Nil {
		return entity.Friendship{}, definitions.ErrInvalidUserInput
	}
	if requesterID == addresseeID {
		return entity.Friendship{}, definitions.ErrInvalidUserInput
	}

	existing, err := uc.repo.GetByUsers(ctx, requesterID, addresseeID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return entity.Friendship{}, err
	}
	if existing != nil {
		return entity.Friendship{}, definitions.ErrFriendshipAlreadyExists
	}

	f := entity.NewFriendship(requesterID, addresseeID)
	return uc.repo.Create(ctx, f)
}

func (uc *UseCase) AcceptRequest(ctx context.Context, friendshipID, userID uuid.UUID) error {
	f, err := uc.repo.GetByID(ctx, friendshipID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return definitions.ErrFriendshipNotFound
		}
		return err
	}

	if f.AddresseeID != userID {
		return definitions.ErrForbidden
	}
	if f.Status != definitions.FriendshipStatusPending {
		return definitions.ErrFriendshipInvalidStatus
	}

	return uc.repo.UpdateStatus(ctx, friendshipID, definitions.FriendshipStatusAccepted)
}

func (uc *UseCase) DeclineRequest(ctx context.Context, friendshipID, userID uuid.UUID) error {
	f, err := uc.repo.GetByID(ctx, friendshipID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return definitions.ErrFriendshipNotFound
		}
		return err
	}

	if f.AddresseeID != userID {
		return definitions.ErrForbidden
	}
	if f.Status != definitions.FriendshipStatusPending {
		return definitions.ErrFriendshipInvalidStatus
	}

	return uc.repo.UpdateStatus(ctx, friendshipID, definitions.FriendshipStatusDeclined)
}

func (uc *UseCase) RemoveFriend(ctx context.Context, friendshipID, userID uuid.UUID) error {
	f, err := uc.repo.GetByID(ctx, friendshipID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return definitions.ErrFriendshipNotFound
		}
		return err
	}

	if f.RequesterID != userID && f.AddresseeID != userID {
		return definitions.ErrForbidden
	}

	return uc.repo.Delete(ctx, friendshipID)
}

func (uc *UseCase) GetFriends(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error) {
	if userID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetFriends(ctx, userID)
}

func (uc *UseCase) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error) {
	if userID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}
	return uc.repo.GetPendingRequests(ctx, userID)
}

func (uc *UseCase) AreFriends(ctx context.Context, userA, userB uuid.UUID) (bool, error) {
	f, err := uc.repo.GetByUsers(ctx, userA, userB)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return f != nil && f.Status == definitions.FriendshipStatusAccepted, nil
}
