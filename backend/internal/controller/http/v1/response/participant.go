package response

type ParticipantResponse struct {
	ID      string `json:"id"`
	EventID string `json:"event_id"`
	UserID  string `json:"user_id"`
}
