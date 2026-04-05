package entity

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID         `db:"id"`
	UserID    uuid.UUID         `db:"user_id"`
	Type      string            `db:"type"`
	Payload   map[string]string `db:"payload"`
	IsRead    bool              `db:"is_read"`
	CreatedAt time.Time         `db:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType string, payload map[string]string) Notification {
	return Notification{
		UserID:  userID,
		Type:    notifType,
		Payload: payload,
	}
}
