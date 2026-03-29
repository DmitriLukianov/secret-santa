package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventStatusDraft              EventStatus = "draft"               // Черновик (только организатор может редактировать)
	EventStatusInvitationOpen     EventStatus = "invitation_open"     // Открыт набор участников по ссылке
	EventStatusRegistrationClosed EventStatus = "registration_closed" // Набор участников закрыт
	EventStatusDrawingPending     EventStatus = "drawing_pending"     // Готов к жеребьёвке
	EventStatusDrawingDone        EventStatus = "drawing_done"        // Жеребьёвка проведена, назначения видны
	EventStatusActive             EventStatus = "active"              // Игра идёт (можно отмечать отправку подарков)
	EventStatusFinished           EventStatus = "finished"            // Событие завершено
	EventStatusCancelled          EventStatus = "cancelled"           // Отменено
)

var ErrInvalidEventState = errors.New("invalid event state transition")

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
		Status:          EventStatusDraft,
		MaxParticipants: maxParticipants,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// CanTransitionTo — проверяет возможность перехода
func (e Event) CanTransitionTo(newStatus EventStatus) bool {
	switch e.Status {
	case EventStatusDraft:
		return newStatus == EventStatusInvitationOpen || newStatus == EventStatusCancelled
	case EventStatusInvitationOpen:
		return newStatus == EventStatusRegistrationClosed || newStatus == EventStatusCancelled
	case EventStatusRegistrationClosed:
		return newStatus == EventStatusDrawingPending || newStatus == EventStatusCancelled
	case EventStatusDrawingPending:
		return newStatus == EventStatusDrawingDone
	case EventStatusDrawingDone:
		return newStatus == EventStatusActive || newStatus == EventStatusCancelled
	case EventStatusActive:
		return newStatus == EventStatusFinished || newStatus == EventStatusCancelled
	default:
		return false
	}
}

// TransitionTo — выполняет переход статуса с проверкой
func (e *Event) TransitionTo(newStatus EventStatus) error {
	if !e.CanTransitionTo(newStatus) {
		return ErrInvalidEventState
	}
	e.Status = newStatus
	e.UpdatedAt = time.Now()
	return nil
}

// ====================== Вспомогательные методы ======================

func (e Event) IsDrawable() bool {
	return e.Status == EventStatusRegistrationClosed || e.Status == EventStatusDrawingPending
}

func (e Event) CanAddParticipants() bool {
	return e.Status == EventStatusDraft || e.Status == EventStatusInvitationOpen
}

func (e Event) CanEdit() bool {
	return e.Status != EventStatusFinished && e.Status != EventStatusCancelled
}

func (e Event) IsActive() bool {
	return e.Status == EventStatusActive || e.Status == EventStatusDrawingDone
}
