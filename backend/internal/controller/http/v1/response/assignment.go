package response

type AssignmentResponse struct {
	ID         string `json:"id"`
	EventID    string `json:"event_id"`
	GiverID    string `json:"giver_id"`
	ReceiverID string `json:"receiver_id"`
}
