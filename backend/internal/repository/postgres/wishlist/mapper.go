package wishlist

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanWishlist(row pgx.Row) (*entity.Wishlist, error) {
	var w entity.Wishlist
	err := row.Scan(
		&w.ID,
		&w.ParticipantID,
		&w.Visibility,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func scanWishlistItems(rows pgx.Rows) ([]entity.WishlistItem, error) {
	var items []entity.WishlistItem
	for rows.Next() {
		var item entity.WishlistItem
		err := rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Title,
			&item.Link,
			&item.ImageURL,
			&item.Comment,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func scanWishlistItem(row pgx.Row) (*entity.WishlistItem, error) {
	var item entity.WishlistItem
	err := row.Scan(
		&item.ID,
		&item.WishlistID,
		&item.Title,
		&item.Link,
		&item.ImageURL,
		&item.Comment,
		&item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}
