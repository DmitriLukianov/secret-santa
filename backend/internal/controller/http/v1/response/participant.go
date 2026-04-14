package response

import "time"

type ParticipantResponse struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventId"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UserName  string    `json:"userName"`
	UserEmail string    `json:"userEmail"`
}
