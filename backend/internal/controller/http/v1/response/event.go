package response

import (
	"time"
)

type EventResponse struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	OrganizerNotes *string    `json:"organizerNotes,omitempty"`
	OrganizerID    string     `json:"organizerId"`
	StartDate      time.Time  `json:"startDate"`
	DrawDate       *time.Time `json:"drawDate,omitempty"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}
