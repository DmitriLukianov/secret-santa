package entity

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID            uuid.UUID  `db:"id"`
	ParticipantID *uuid.UUID `db:"participant_id"` // nil для персонального вишлиста
	UserID        *uuid.UUID `db:"user_id"`         // nil для вишлиста участника события
	Visibility    string     `db:"visibility"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

type WishlistItem struct {
	ID         uuid.UUID `db:"id"`
	WishlistID uuid.UUID `db:"wishlist_id"`
	Title      string    `db:"title"`
	Link       *string   `db:"link"`
	ImageURL   *string   `db:"image_url"`
	Price      *float64  `db:"price"`
	CreatedAt  time.Time `db:"created_at"`
	Wishlist   *Wishlist `db:"-"`
}

func NewWishlist(participantID uuid.UUID, visibility string) Wishlist {
	return Wishlist{
		ParticipantID: &participantID,
		Visibility:    visibility,
	}
}

func NewPersonalWishlist(userID uuid.UUID) Wishlist {
	return Wishlist{
		UserID:     &userID,
		Visibility: "public",
	}
}

func NewWishlistItem(wishlistID uuid.UUID, title string, link, imageURL *string, price *float64) WishlistItem {
	return WishlistItem{
		WishlistID: wishlistID,
		Title:      title,
		Link:       link,
		ImageURL:   imageURL,
		Price:      price,
	}
}
