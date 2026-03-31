package request

type CreateUserRequest struct {
	Name          string `json:"name" validate:"required,min=2"`
	Email         string `json:"email" validate:"required,email"`
	OAuthID       string `json:"oauthId" validate:"required"`
	OAuthProvider string `json:"oauthProvider" validate:"required,oneof=github vk google"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=2"`
	Email *string `json:"email" validate:"omitempty,email"`
}
