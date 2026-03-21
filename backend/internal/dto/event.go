package dto

type CreateEventInput struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	OrganizerID string `json:"organizer_id" validate:"required"`
	StartDate   string `json:"start_date"`
	DrawDate    string `json:"draw_date"`
	EndDate     string `json:"end_date"`
}

type UpdateEventInput struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
