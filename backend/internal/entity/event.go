package entity

import (
	"time"

	"github.com/google/uuid"

	"secret-santa-backend/internal/definitions"
)

type Event struct {
	ID             uuid.UUID               `db:"id"`
	Title          string                  `db:"title"`
	OrganizerNotes *string                 `db:"organizer_notes"`
	OrganizerID    uuid.UUID               `db:"organizer_id"`
	StartDate      time.Time               `db:"start_date"`
	DrawDate       *time.Time              `db:"draw_date"`
	Budget         *int                    `db:"budget"`
	Status         definitions.EventStatus `db:"status"`
	CreatedAt      time.Time               `db:"created_at"`
	UpdatedAt      time.Time               `db:"updated_at"`
}

func NewEvent(
	title string,
	organizerID uuid.UUID,
	organizerNotes *string,
	startDate time.Time,
	drawDate *time.Time,
	budget *int,
) Event {
	return Event{
		Title:          title,
		OrganizerNotes: organizerNotes,
		OrganizerID:    organizerID,
		StartDate:      startDate,
		DrawDate:       drawDate,
		Budget:         budget,
		Status:         definitions.EventStatusRegistration,
	}
}

func (e *Event) TransitionTo(newStatus definitions.EventStatus) error {
	switch {
	case e.Status == definitions.EventStatusRegistration && newStatus == definitions.EventStatusGifting:
	case e.Status == definitions.EventStatusGifting && newStatus == definitions.EventStatusFinished:
	default:
		return definitions.ErrInvalidEventState
	}
	e.Status = newStatus
	e.UpdatedAt = time.Now()
	return nil
}

func (e Event) IsDrawable() bool {
	return e.Status == definitions.EventStatusRegistration
}

func (e Event) CanAddParticipants() bool {
	return e.Status == definitions.EventStatusRegistration
}

func (e Event) CanEdit() bool {
	return e.Status == definitions.EventStatusRegistration
}

func (e Event) IsActive() bool {
	return e.Status == definitions.EventStatusGifting
}

func (e Event) IsFinished() bool {
	return e.Status == definitions.EventStatusFinished
}
