package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateEventInput — данные для создания нового события
type CreateEventInput struct {
	Name            string     `json:"name" validate:"required,min=3,max=100"`
	Description     *string    `json:"description" validate:"omitempty,max=1000"`
	Rules           *string    `json:"rules" validate:"omitempty,max=2000"`
	Recommendations *string    `json:"recommendations" validate:"omitempty,max=2000"`
	StartDate       *time.Time `json:"startDate" validate:"omitempty"`
	DrawDate        *time.Time `json:"drawDate" validate:"omitempty"`
	EndDate         *time.Time `json:"endDate" validate:"omitempty"`
	MaxParticipants int        `json:"maxParticipants" validate:"min=0"`
}

// UpdateEventInput — данные для обновления события (partial update)
type UpdateEventInput struct {
	Name            *string    `json:"name" validate:"omitempty,min=3,max=100"`
	Description     *string    `json:"description" validate:"omitempty,max=1000"`
	Rules           *string    `json:"rules" validate:"omitempty,max=2000"`
	Recommendations *string    `json:"recommendations" validate:"omitempty,max=2000"`
	StartDate       *time.Time `json:"startDate" validate:"omitempty"`
	DrawDate        *time.Time `json:"drawDate" validate:"omitempty"`
	EndDate         *time.Time `json:"endDate" validate:"omitempty"`
	Status          *string    `json:"status" validate:"omitempty,oneof=draft active finished cancelled"`
	MaxParticipants *int       `json:"maxParticipants" validate:"omitempty,min=0"`
}

// EventResponse — ответ с полным событием (для фронтенда)
type EventResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Rules           *string    `json:"rules,omitempty"`
	Recommendations *string    `json:"recommendations,omitempty"`
	OrganizerID     uuid.UUID  `json:"organizerId"`
	StartDate       *time.Time `json:"startDate,omitempty"`
	DrawDate        *time.Time `json:"drawDate,omitempty"`
	EndDate         *time.Time `json:"endDate,omitempty"`
	Status          string     `json:"status"`
	MaxParticipants int        `json:"maxParticipants"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
