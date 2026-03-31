package entity

import (
	"time"

	"github.com/google/uuid"
)

type Assignment struct {
	ID         uuid.UUID `db:"id"`
	EventID    uuid.UUID `db:"event_id"`
	GiverID    uuid.UUID `db:"giver_id"`
	ReceiverID uuid.UUID `db:"receiver_id"`
	CreatedAt  time.Time `db:"created_at"`
}

func NewAssignment(eventID, giverID, receiverID uuid.UUID) Assignment {
	return Assignment{
		ID:         uuid.New(),
		EventID:    eventID,
		GiverID:    giverID,
		ReceiverID: receiverID,
		CreatedAt:  time.Now(),
	}
}
