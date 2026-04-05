package entity

import (
	"time"

	"github.com/google/uuid"
)

type Friendship struct {
	ID          uuid.UUID `db:"id"`
	RequesterID uuid.UUID `db:"requester_id"`
	AddresseeID uuid.UUID `db:"addressee_id"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func NewFriendship(requesterID, addresseeID uuid.UUID) Friendship {
	return Friendship{
		RequesterID: requesterID,
		AddresseeID: addresseeID,
		Status:      "pending",
	}
}
