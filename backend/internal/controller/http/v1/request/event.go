package request

import "time"

type CreateEventRequest struct {
	Title           string     `json:"title" validate:"required,min=3"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	OrganizerNotes  *string    `json:"organizerNotes"`
	StartDate       *time.Time `json:"startDate"`
	DrawDate        *time.Time `json:"drawDate"`
	EndDate         *time.Time `json:"endDate"`
	MaxParticipants *int       `json:"maxParticipants"`
	WantParticipate *bool      `json:"wantParticipate"`
}

type UpdateEventRequest struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	Rules           *string    `json:"rules"`
	Recommendations *string    `json:"recommendations"`
	OrganizerNotes  *string    `json:"organizerNotes"`
	StartDate       *time.Time `json:"startDate"`
	DrawDate        *time.Time `json:"drawDate"`
	ClearDrawDate   bool       `json:"clearDrawDate"`
	EndDate         *time.Time `json:"endDate"`
	MaxParticipants *int       `json:"maxParticipants"`
}
