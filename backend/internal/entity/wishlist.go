package entity

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID            uuid.UUID `db:"id"`
	ParticipantID uuid.UUID `db:"participant_id"`
	Visibility    string    `db:"visibility"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type WishlistItem struct {
	ID         uuid.UUID `db:"id"`
	WishlistID uuid.UUID `db:"wishlist_id"`
	Title      string    `db:"title"`
	Link       *string   `db:"link"`
	ImageURL   *string   `db:"image_url"`
	Comment    *string   `db:"comment"`
	Price      *float64  `db:"price"`
	CreatedAt  time.Time `db:"created_at"`
	Wishlist   *Wishlist `db:"-"`
}

func NewWishlist(participantID uuid.UUID, visibility string) Wishlist {
	return Wishlist{
		ParticipantID: participantID,
		Visibility:    visibility,
	}
}

func NewWishlistItem(wishlistID uuid.UUID, title string, link, imageURL, comment *string, price *float64) WishlistItem {
	return WishlistItem{
		WishlistID: wishlistID,
		Title:      title,
		Link:       link,
		ImageURL:   imageURL,
		Comment:    comment,
		Price:      price,
	}
}
