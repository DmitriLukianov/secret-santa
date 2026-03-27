package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID         uuid.UUID  `db:"id"`
	EventID    uuid.UUID  `db:"event_id"`
	UserID     uuid.UUID  `db:"user_id"`
	Role       string     `db:"role"`
	GiftSent   bool       `db:"gift_sent"`
	GiftSentAt *time.Time `db:"gift_sent_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func NewParticipant(eventID, userID uuid.UUID, role string) Participant {
	now := time.Now()
	return Participant{
		ID:        uuid.New(),
		EventID:   eventID,
		UserID:    userID,
		Role:      role,
		GiftSent:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

const (
	ParticipantRoleOrganizer   = "organizer"
	ParticipantRoleParticipant = "participant"
)
