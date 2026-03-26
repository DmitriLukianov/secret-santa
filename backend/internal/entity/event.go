package entity

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusActive    EventStatus = "active"
	EventStatusFinished  EventStatus = "finished"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID              uuid.UUID   `db:"id"`
	Title           string      `db:"title"`
	Description     *string     `db:"description"`
	Rules           *string     `db:"rules"`
	Recommendations *string     `db:"recommendations"`
	OrganizerID     uuid.UUID   `db:"organizer_id"`
	StartDate       time.Time   `db:"start_date"`
	DrawDate        *time.Time  `db:"draw_date"`
	EndDate         time.Time   `db:"end_date"`
	Status          EventStatus `db:"status"`
	MaxParticipants int         `db:"max_participants"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
}

func NewEvent(
	title string,
	organizerID uuid.UUID,
	description, rules, recommendations *string,
	startDate, drawDate, endDate time.Time,
	maxParticipants int,
) Event {
	now := time.Now()

	return Event{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		Rules:           rules,
		Recommendations: recommendations,
		OrganizerID:     organizerID,
		StartDate:       startDate,
		DrawDate:        &drawDate,
		EndDate:         endDate,
		Status:          EventStatusDraft,
		MaxParticipants: maxParticipants,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func (e Event) IsDrawable() bool {
	return e.Status == EventStatusDraft
}

func (e *Event) MarkAsDrawn() {
	if e.Status == EventStatusDraft {
		e.Status = EventStatusActive
		e.UpdatedAt = time.Now()
	}
}

func (e Event) CanBeFinished() bool {
	return e.Status == EventStatusDraft || e.Status == EventStatusActive
}

func (e *Event) MarkAsFinished() {
	e.Status = EventStatusFinished
	e.UpdatedAt = time.Now()
}
