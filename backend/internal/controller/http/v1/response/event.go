package response

import "time"

type EventResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Description     *string    `json:"description,omitempty"`
	Rules           *string    `json:"rules,omitempty"`
	Recommendations *string    `json:"recommendations,omitempty"`
	OrganizerID     string     `json:"organizer_id"`
	StartDate       *time.Time `json:"start_date,omitempty"`
	DrawDate        *time.Time `json:"draw_date,omitempty"`
	EndDate         *time.Time `json:"end_date,omitempty"`
	Status          string     `json:"status"`
	MaxParticipants int        `json:"max_participants"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
