package request

import "time"

type CreateEventRequest struct {
	Name            string     `json:"name"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	StartDate       *time.Time `json:"start_date"`
	DrawDate        *time.Time `json:"draw_date"`
	EndDate         *time.Time `json:"end_date"`
	MaxParticipants int        `json:"max_participants"`
}

type UpdateEventRequest struct {
	Name            *string    `json:"name"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	StartDate       *time.Time `json:"start_date"`
	DrawDate        *time.Time `json:"draw_date"`
	EndDate         *time.Time `json:"end_date"`
	Status          *string    `json:"status"`
	MaxParticipants *int       `json:"max_participants"`
}
