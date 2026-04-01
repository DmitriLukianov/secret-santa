package request

type CreateWishlistRequest struct {
	EventID    string `json:"eventId" validate:"required,uuid"`
	Visibility string `json:"visibility" validate:"required,oneof=public santa_only"`
}

type CreateWishlistItemRequest struct {
	Title    string `json:"title" validate:"required"`
	Link     string `json:"link,omitempty"`
	ImageURL string `json:"imageURL,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// NEW: запрос на обновление товара (все поля кроме title — опциональные)
type UpdateWishlistItemRequest struct {
	Title    string `json:"title" validate:"required"`
	Link     string `json:"link,omitempty"`
	ImageURL string `json:"imageURL,omitempty"`
	Comment  string `json:"comment,omitempty"`
}
