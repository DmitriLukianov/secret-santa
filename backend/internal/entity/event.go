package entity

import "time"

type Event struct {
	ID          string
	Name        string
	Description string
	OrganizerID string
	StartDate   time.Time
	DrawDate    time.Time
	EndDate     time.Time
	CreatedAt   time.Time
}
