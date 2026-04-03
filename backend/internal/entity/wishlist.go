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
	CreatedAt  time.Time `db:"created_at"`
	Wishlist   *Wishlist `db:"-"`
}

// NewWishlist — чистый конструктор, ID и таймстемпы даёт БД
func NewWishlist(participantID uuid.UUID, visibility string) Wishlist {
	return Wishlist{
		ParticipantID: participantID,
		Visibility:    visibility,
		// ID, CreatedAt, UpdatedAt будут заполнены БД
	}
}

// NewWishlistItem — чистый конструктор
func NewWishlistItem(wishlistID uuid.UUID, title string, link, imageURL, comment *string) WishlistItem {
	return WishlistItem{
		WishlistID: wishlistID,
		Title:      title,
		Link:       link,
		ImageURL:   imageURL,
		Comment:    comment,
		// ID и CreatedAt будут заполнены БД
	}
}
