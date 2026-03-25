package entity

import (
	"time"

	"github.com/google/uuid"
)

// Event — событие «Тайный Санта»
type Event struct {
	ID              uuid.UUID  `db:"id"`
	Name            string     `db:"name"`
	Description     *string    `db:"description"` // может быть null
	Rules           *string    `db:"rules"`
	Recommendations *string    `db:"recommendations"`
	OrganizerID     uuid.UUID  `db:"organizer_id"`
	StartDate       *time.Time `db:"start_date"`
	DrawDate        *time.Time `db:"draw_date"` // дата жеребьёвки
	EndDate         *time.Time `db:"end_date"`
	Status          string     `db:"status"` // draft | active | finished | cancelled
	MaxParticipants int        `db:"max_participants"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

// NewEvent — конструктор для создания нового события
func NewEvent(
	name string,
	organizerID uuid.UUID,
	description, rules, recommendations *string,
	startDate, drawDate, endDate *time.Time,
	maxParticipants int,
) Event {
	now := time.Now()

	return Event{
		ID:              uuid.New(),
		Name:            name,
		Description:     description,
		Rules:           rules,
		Recommendations: recommendations,
		OrganizerID:     organizerID,
		StartDate:       startDate,
		DrawDate:        drawDate,
		EndDate:         endDate,
		Status:          "draft", // по умолчанию черновик
		MaxParticipants: maxParticipants,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
