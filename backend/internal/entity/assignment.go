package entity

import (
	"time"

	"github.com/google/uuid"
)

// Assignment — назначение «кто кому дарит» в событии
type Assignment struct {
	ID         uuid.UUID `db:"id"`
	EventID    uuid.UUID `db:"event_id"`
	GiverID    uuid.UUID `db:"giver_id"`    // кто дарит
	ReceiverID uuid.UUID `db:"receiver_id"` // кому дарит
	CreatedAt  time.Time `db:"created_at"`
}

// NewAssignment — конструктор
func NewAssignment(eventID, giverID, receiverID uuid.UUID) Assignment {
	return Assignment{
		ID:         uuid.New(),
		EventID:    eventID,
		GiverID:    giverID,
		ReceiverID: receiverID,
		CreatedAt:  time.Now(),
	}
}
