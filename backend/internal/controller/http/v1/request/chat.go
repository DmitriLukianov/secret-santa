package request

type SendMessageRequest struct {
	Content string `json:"content" validate:"required"`
}
