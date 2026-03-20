package domain

import "time"

type Participant struct {
	ID       string    `json:"id"`
	EventID  string    `json:"event_id"`
	UserID   string    `json:"user_id"`
	Status   string    `json:"status"`
	JoinedAt time.Time `json:"joined_at"`
}
