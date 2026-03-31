package response

import "time"

type AssignmentResponse struct {
	ID         string    `json:"id"`
	EventID    string    `json:"eventId"`
	GiverID    string    `json:"giverId"`
	ReceiverID string    `json:"receiverId"`
	CreatedAt  time.Time `json:"createdAt"`
}
