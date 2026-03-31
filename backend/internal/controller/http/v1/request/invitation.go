package request

import "github.com/google/uuid"

type CreateInvitationRequest struct {
	EventID   uuid.UUID `json:"eventId" validate:"required"`
	ExpiresIn string    `json:"expiresIn,omitempty"` // ← теперь string
}

type JoinByInvitationRequest struct {
	Token string `json:"token" validate:"required"`
}
