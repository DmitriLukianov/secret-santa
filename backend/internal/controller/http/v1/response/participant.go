package response

import "time"

type ParticipantResponse struct {
	ID         string     `json:"id"`
	EventID    string     `json:"eventId"`
	UserID     string     `json:"userId"`
	Role       string     `json:"role"`
	GiftSent   bool       `json:"giftSent"`
	GiftSentAt *time.Time `json:"giftSentAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
}
