package dto

// CreateUserInput — входные данные для создания пользователя (из OAuth)
type CreateUserInput struct {
	Name          string `json:"name" validate:"required,min=2"`
	Email         string `json:"email" validate:"required,email"`
	OAuthID       string `json:"oauthId" validate:"required"`
	OAuthProvider string `json:"oauthProvider" validate:"required,oneof=github vk google"`
}

// UpdateUserInput — входные данные для обновления профиля
type UpdateUserInput struct {
	Name  *string `json:"name" validate:"omitempty,min=2"`
	Email *string `json:"email" validate:"omitempty,email"`
}
