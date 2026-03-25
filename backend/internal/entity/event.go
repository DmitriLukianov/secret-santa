package entity

import (
	"time"

	"github.com/google/uuid"
)

// Event — событие «Тайный Санта»
type Event struct {
	ID              uuid.UUID  `db:"id"`
	Title           string     `db:"title"` // было Name → теперь Title
	Description     *string    `db:"description"`
	Rules           *string    `db:"rules"`
	Recommendations *string    `db:"recommendations"`
	OrganizerID     uuid.UUID  `db:"organizer_id"`
	StartDate       time.Time  `db:"start_date"` // теперь обязательно
	DrawDate        *time.Time `db:"draw_date"`
	EndDate         time.Time  `db:"end_date"` // теперь обязательно
	Status          string     `db:"status"`
	MaxParticipants int        `db:"max_participants"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

// Константы статусов события (можно позже вынести в definitions/constants.go)
const (
	EventStatusDraft     = "draft"
	EventStatusDrawing   = "drawing"
	EventStatusActive    = "active"
	EventStatusFinished  = "finished"
	EventStatusCancelled = "cancelled"
)

// NewEvent — конструктор
func NewEvent(
	title string,
	organizerID uuid.UUID,
	description, rules, recommendations *string,
	startDate, drawDate, endDate *time.Time,
	maxParticipants int,
) Event {
	now := time.Now()

	// Если даты не передали — ставим разумные значения
	if startDate == nil {
		startDate = &now
	}
	if endDate == nil {
		endDate = &now
	}

	return Event{
		ID:              uuid.New(),
		Title:           title,
		Description:     description,
		Rules:           rules,
		Recommendations: recommendations,
		OrganizerID:     organizerID,
		StartDate:       *startDate,
		DrawDate:        drawDate,
		EndDate:         *endDate,
		Status:          EventStatusDraft,
		MaxParticipants: maxParticipants,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
