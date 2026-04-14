package notification

import (
	"context"
	"log/slog"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

const (
	TypeDrawDone        = "draw_done"
	TypeInvitationJoined = "invitation_joined"
	TypeGiftSent        = "gift_sent"
	TypeFriendRequest   = "friend_request"
	TypeFriendAccepted  = "friend_accepted"
)

type UseCase struct {
	repo Repository
	log  *slog.Logger
}

func New(repo Repository) *UseCase {
	return &UseCase{repo: repo}
}

func NewWithLogger(repo Repository, log *slog.Logger) *UseCase {
	return &UseCase{repo: repo, log: log}
}

func (uc *UseCase) Notify(ctx context.Context, userID uuid.UUID, notifType string, payload map[string]string) error {
	n := entity.NewNotification(userID, notifType, payload)
	_, err := uc.repo.Create(ctx, n)
	if err != nil && uc.log != nil {
		uc.log.Error("failed to create notification",
			slog.String("user_id", userID.String()),
			slog.String("type", notifType),
			slog.String("error", err.Error()),
		)
	}
	return err
}

