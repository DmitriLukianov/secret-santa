package request

import "github.com/google/uuid"

type CreateInvitationRequest struct {
	EventID   uuid.UUID `json:"eventId" validate:"required"`
	ExpiresIn string    `json:"expiresIn,omitempty"`
}

type JoinByInvitationRequest struct {
	InvitationLink string `json:"invitationLink" validate:"required"`
}

type SendEmailInvitationRequest struct {
	EventID   string `json:"eventId" validate:"required,uuid"`
	Email     string `json:"email" validate:"required,email"`
	ExpiresIn string `json:"expiresIn,omitempty"`
}
