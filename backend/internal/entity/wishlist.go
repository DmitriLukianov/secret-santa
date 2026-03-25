package entity

import (
	"time"

	"github.com/google/uuid"
)

// Wishlist — список желаемых подарков участника события
type Wishlist struct {
	ID            uuid.UUID `db:"id"`
	ParticipantID uuid.UUID `db:"participant_id"` // привязка к участнику, а не к user
	Visibility    string    `db:"visibility"`     // public | friends | santa_only
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// WishlistItem — отдельный элемент вишлиста (можно добавлять несколько)
type WishlistItem struct {
	ID         uuid.UUID `db:"id"`
	WishlistID uuid.UUID `db:"wishlist_id"`
	Title      string    `db:"title"`
	Link       *string   `db:"link"`
	ImageURL   *string   `db:"image_url"`
	Comment    *string   `db:"comment"`
	CreatedAt  time.Time `db:"created_at"`
}

// Константы видимости (по ТЗ)
const (
	WishlistVisibilityPublic    = "public"     // для всех
	WishlistVisibilityFriends   = "friends"    // только друзья
	WishlistVisibilitySantaOnly = "santa_only" // только Тайный Санта
)

// NewWishlist — конструктор
func NewWishlist(participantID uuid.UUID, visibility string) Wishlist {
	now := time.Now()
	return Wishlist{
		ID:            uuid.New(),
		ParticipantID: participantID,
		Visibility:    visibility,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewWishlistItem — конструктор элемента
func NewWishlistItem(wishlistID uuid.UUID, title string, link, imageURL, comment *string) WishlistItem {
	return WishlistItem{
		ID:         uuid.New(),
		WishlistID: wishlistID,
		Title:      title,
		Link:       link,
		ImageURL:   imageURL,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}
}
