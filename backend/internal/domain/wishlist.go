package domain

import "time"

type Wishlist struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}
