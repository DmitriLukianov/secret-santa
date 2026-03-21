package entity

import "time"

type Wishlist struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Link        string
	ImageURL    string
	Visibility  string
	CreatedAt   time.Time
}
