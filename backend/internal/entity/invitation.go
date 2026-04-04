package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID        uuid.UUID `db:"id"`
	EventID   uuid.UUID `db:"event_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewInvitation(eventID, createdBy uuid.UUID, expiresIn time.Duration) Invitation {
	if expiresIn == 0 {
		expiresIn = 7 * 24 * time.Hour
	}

	return Invitation{
		EventID:   eventID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(expiresIn),
		CreatedBy: createdBy,
	}
}

func (i Invitation) IsValid() bool {
	return time.Now().Before(i.ExpiresAt)
}
