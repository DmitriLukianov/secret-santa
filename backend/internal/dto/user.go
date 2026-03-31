package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserInput struct {
	Name          string `json:"name" validate:"required,min=2"`
	Email         string `json:"email" validate:"required,email"`
	OAuthID       string `json:"oauthId" validate:"required"`
	OAuthProvider string `json:"oauthProvider" validate:"required,oneof=vk google github"`
}

type UpdateUserInput struct {
	Name  *string `json:"name" validate:"omitempty,min=2"`
	Email *string `json:"email" validate:"omitempty,email"`
}

type UserResponse struct { // новый — для ответов API
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	OAuthProvider string    `json:"oauthProvider"`
	CreatedAt     time.Time `json:"createdAt"`
}
