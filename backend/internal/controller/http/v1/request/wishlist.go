package request

type WishlistRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	ImageURL    string `json:"image_url"`
	Visibility  string `json:"visibility"`
}
