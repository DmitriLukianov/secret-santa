package dto

type CreateUserInput struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
