package dto

type CreateWishlistInput struct {
	UserID      string
	Title       string
	Description string
	Link        string
	ImageURL    string
	Visibility  string
}

type UpdateWishlistInput struct {
	Title       *string
	Description *string
	Link        *string
	ImageURL    *string
	Visibility  *string
}
