package entity

import (
	"time"

	"github.com/google/uuid"
)

// Participant — участник события «Тайный Санта»
type Participant struct {
	ID         uuid.UUID  `db:"id"`
	EventID    uuid.UUID  `db:"event_id"`
	UserID     uuid.UUID  `db:"user_id"`
	Role       string     `db:"role"` // "organizer" | "participant"
	GiftSent   bool       `db:"gift_sent"`
	GiftSentAt *time.Time `db:"gift_sent_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

// NewParticipant — конструктор
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

// Константы ролей (можно позже вынести в definitions/constants.go)
const (
	ParticipantRoleOrganizer   = "organizer"
	ParticipantRoleParticipant = "participant"
)
