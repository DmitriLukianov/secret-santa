package request

type CreateWishlistRequest struct {
	Visibility string `json:"visibility" validate:"required,oneof=public friends santa_only"`
}

type CreateWishlistItemRequest struct {
	Title    string  `json:"title" validate:"required,min=2"`
	Link     *string `json:"link"`
	ImageURL *string `json:"imageUrl"`
	Comment  *string `json:"comment"`
}
