package dto

type SendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}
