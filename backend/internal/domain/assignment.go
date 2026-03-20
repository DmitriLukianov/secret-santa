package domain

import "time"

type Assignment struct {
	ID         string    `json:"id"`
	EventID    string    `json:"event_id"`
	GiverID    string    `json:"giver_id"`
	ReceiverID string    `json:"receiver_id"`
	CreatedAt  time.Time `json:"created_at"`
}
