package dto

import "time"

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
