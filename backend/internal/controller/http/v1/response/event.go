package response

import "time"

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
