package response

import (
	"secret-santa-backend/internal/entity"
	"time"
)

type EventResponse struct {
	ID              string     `json:"id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	Rules           *string    `json:"rules,omitempty"`
	Recommendations *string    `json:"recommendations,omitempty"`
	OrganizerID     string     `json:"organizerId"`
	StartDate       time.Time  `json:"startDate"`
	DrawDate        *time.Time `json:"drawDate,omitempty"`
	EndDate         time.Time  `json:"endDate"`
	Status          string     `json:"status"`
	MaxParticipants int        `json:"maxParticipants"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func EventResponseFromEntity(e entity.Event) EventResponse {
	return EventResponse{
		ID:              e.ID.String(),
		Title:           e.Title,
		Description:     e.Description,
		Rules:           e.Rules,
		Recommendations: e.Recommendations,
		OrganizerID:     e.OrganizerID.String(),
		StartDate:       e.StartDate,
		DrawDate:        e.DrawDate,
		EndDate:         e.EndDate,
		Status:          string(e.Status),
		MaxParticipants: e.MaxParticipants,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
