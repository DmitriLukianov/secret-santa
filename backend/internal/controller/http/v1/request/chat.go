package request

type SendMessageRequest struct {
	Text string `json:"text" validate:"required"`
}
