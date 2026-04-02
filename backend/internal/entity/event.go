package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"secret-santa-backend/internal/definitions" // ← добавлен импорт
)

type Event struct {
	ID              uuid.UUID               `db:"id"`
	Title           string                  `db:"title"`
	Description     *string                 `db:"description"`
	Rules           *string                 `db:"rules"`
	Recommendations *string                 `db:"recommendations"`
	OrganizerID     uuid.UUID               `db:"organizer_id"`
	StartDate       time.Time               `db:"start_date"`
	DrawDate        *time.Time              `db:"draw_date"`
	EndDate         time.Time               `db:"end_date"`
	Status          definitions.EventStatus `db:"status"`
	MaxParticipants int                     `db:"max_participants"`
	CreatedAt       time.Time               `db:"created_at"`
	UpdatedAt       time.Time               `db:"updated_at"`
}

// NewEvent создаёт событие в статусе draft
func NewEvent(
	title string,
	organizerID uuid.UUID,
	description, rules, recommendations *string,
	startDate, endDate time.Time,
	drawDate *time.Time,
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
		DrawDate:        drawDate,
		EndDate:         endDate,
		Status:          definitions.EventStatusDraft,
		MaxParticipants: maxParticipants,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// CanTransitionTo — проверяет возможность перехода
func (e Event) CanTransitionTo(newStatus definitions.EventStatus) bool {
	switch e.Status {
	case definitions.EventStatusDraft:
		return newStatus == definitions.EventStatusInvitationOpen || newStatus == definitions.EventStatusCancelled
	case definitions.EventStatusInvitationOpen:
		return newStatus == definitions.EventStatusRegistrationClosed || newStatus == definitions.EventStatusCancelled
	case definitions.EventStatusRegistrationClosed:
		return newStatus == definitions.EventStatusDrawingPending || newStatus == definitions.EventStatusDrawingDone || newStatus == definitions.EventStatusCancelled
	case definitions.EventStatusDrawingPending:
		return newStatus == definitions.EventStatusDrawingDone || newStatus == definitions.EventStatusCancelled
	case definitions.EventStatusDrawingDone:
		return newStatus == definitions.EventStatusActive || newStatus == definitions.EventStatusCancelled
	case definitions.EventStatusActive:
		return newStatus == definitions.EventStatusFinished || newStatus == definitions.EventStatusCancelled
	default:
		return false
	}
}

// TransitionTo — выполняет переход статуса с проверкой
func (e *Event) TransitionTo(newStatus definitions.EventStatus) error {
	if !e.CanTransitionTo(newStatus) {
		return ErrInvalidEventState
	}
	e.Status = newStatus
	e.UpdatedAt = time.Now()
	return nil
}

// ====================== Вспомогательные методы ======================

func (e Event) IsDrawable() bool {
	return e.Status == definitions.EventStatusRegistrationClosed || e.Status == definitions.EventStatusDrawingPending
}

func (e Event) CanAddParticipants() bool {
	return e.Status == definitions.EventStatusDraft || e.Status == definitions.EventStatusInvitationOpen
}

func (e Event) CanEdit() bool {
	return e.Status != definitions.EventStatusFinished && e.Status != definitions.EventStatusCancelled
}

func (e Event) IsActive() bool {
	return e.Status == definitions.EventStatusActive || e.Status == definitions.EventStatusDrawingDone
}

var ErrInvalidEventState = errors.New("invalid event state transition")
