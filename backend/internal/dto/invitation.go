package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateInvitationInput struct {
	EventID   uuid.UUID     `json:"eventId" validate:"required"`
	ExpiresIn time.Duration `json:"expiresIn"`
}

type JoinByInvitationInput struct {
	Token  string    `json:"token" validate:"required"`
	UserID uuid.UUID `json:"-"`
}

type InvitationResponse struct {
	InviteURL string    `json:"inviteUrl"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
