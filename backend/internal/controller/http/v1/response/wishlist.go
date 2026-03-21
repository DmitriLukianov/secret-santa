package response

type WishlistResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	ImageURL    string `json:"image_url"`
	Visibility  string `json:"visibility"`
}
