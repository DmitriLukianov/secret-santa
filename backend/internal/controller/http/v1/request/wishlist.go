package request

type CreateWishlistRequest struct {
	EventID    string `json:"eventId" validate:"required,uuid"`
	Visibility string `json:"visibility" validate:"required,oneof=public santa_only"`
}

type CreateWishlistItemRequest struct {
	Title    string   `json:"title" validate:"required"`
	Link     string   `json:"link,omitempty"`
	ImageURL string   `json:"imageURL,omitempty"`
	Price    *float64 `json:"price,omitempty"`
}

type UpdateWishlistItemRequest struct {
	Title    string   `json:"title" validate:"required"`
	Link     string   `json:"link,omitempty"`
	ImageURL string   `json:"imageURL,omitempty"`
	Price    *float64 `json:"price,omitempty"`
}
