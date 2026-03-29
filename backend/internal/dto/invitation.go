package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateInvitationInput — данные для генерации приглашения (организатор)
type CreateInvitationInput struct {
	EventID   uuid.UUID     `json:"eventId" validate:"required"`
	ExpiresIn time.Duration `json:"expiresIn"` // например 7 * 24 * time.Hour
}

// JoinByInvitationInput — данные при переходе по ссылке
type JoinByInvitationInput struct {
	Token  string    `json:"token" validate:"required"`
	UserID uuid.UUID `json:"-"` // заполняется из middleware
}

// InvitationResponse — ответ при генерации ссылки
type InvitationResponse struct {
	InviteURL string    `json:"inviteUrl"` // полная ссылка для приглашения
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}
