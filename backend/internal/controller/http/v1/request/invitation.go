package request

import "github.com/google/uuid"

type CreateInvitationRequest struct {
	EventID   uuid.UUID `json:"eventId" validate:"required"`
	ExpiresIn string    `json:"expiresIn,omitempty"`
}

type JoinByInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}

type SendEmailInvitationRequest struct {
	EventID        string `json:"eventId" validate:"required,uuid"`
	RecipientEmail string `json:"recipientEmail" validate:"required,email"`
	ExpiresIn      string `json:"expiresIn,omitempty"`
}
