package dto

import "time"

// CreateEventInput — данные для создания события
type CreateEventInput struct {
	Title           string     `json:"title" validate:"required,min=3"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	StartDate       time.Time  `json:"startDate" validate:"required"`
	DrawDate        *time.Time `json:"drawDate"`
	EndDate         time.Time  `json:"endDate" validate:"required"`
	MaxParticipants int        `json:"maxParticipants" validate:"min=2"`
}

// UpdateEventInput — partial update (все поля pointer)
type UpdateEventInput struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	StartDate       *time.Time `json:"startDate"`
	DrawDate        *time.Time `json:"drawDate"`
	EndDate         *time.Time `json:"endDate"`
	MaxParticipants *int       `json:"maxParticipants"`
}

// EventResponse — ответ API
type EventResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Rules           *string    `json:"rules,omitempty"`
	Recommendations *string    `json:"recommendations,omitempty"`
	OrganizerID     string     `json:"organizerId"`
	StartDate       time.Time  `json:"startDate"`
	DrawDate        *time.Time `json:"drawDate,omitempty"`
	EndDate         time.Time  `json:"endDate"`
	Status          string     `json:"status"`
	MaxParticipants int        `json:"maxParticipants"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}
