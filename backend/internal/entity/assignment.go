package entity

import "time"

type Assignment struct {
	ID         string
	EventID    string
	GiverID    string
	ReceiverID string
	CreatedAt  time.Time
}
