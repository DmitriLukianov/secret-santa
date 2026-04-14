package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID        uuid.UUID `db:"id"`
	EventID   uuid.UUID `db:"event_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	// Дополнительные поля из JOIN с users (заполняются только в GetByEvent)
	UserName  string `db:"user_name"`
	UserEmail string `db:"user_email"`
}

func NewParticipant(eventID, userID uuid.UUID) Participant {
	return Participant{
		EventID: eventID,
		UserID:  userID,
	}
}
