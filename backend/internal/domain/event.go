package domain

import "time"

type Event struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OrganizerID string    `json:"organizer_id"`
	StartDate   time.Time `json:"start_date"`
	DrawDate    time.Time `json:"draw_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}
