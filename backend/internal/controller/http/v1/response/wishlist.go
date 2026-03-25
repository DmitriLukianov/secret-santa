package response

import "time"

type WishlistResponse struct {
	ID            string    `json:"id"`
	ParticipantID string    `json:"participantId"`
	Visibility    string    `json:"visibility"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type WishlistItemResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Link      *string   `json:"link"`
	ImageURL  *string   `json:"imageUrl"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
