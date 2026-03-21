package request

type CreateEventRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OrganizerID string `json:"organizer_id"`
	StartDate   string `json:"start_date"`
	DrawDate    string `json:"draw_date"`
	EndDate     string `json:"end_date"`
}

type UpdateEventRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
