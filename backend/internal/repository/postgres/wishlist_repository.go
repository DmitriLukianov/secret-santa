package postgres

import (
	"context"
	"secret-santa-backend/internal/domain"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WishlistRepository struct {
	db *pgxpool.Pool
}

func NewWishlistRepository(db *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, wishlist domain.Wishlist) error {

	query := `
	INSERT INTO wishlists (event_id, user_id, text)
	VALUES ($1,$2,$3)
	`

	_, err := r.db.Exec(ctx, query,
		wishlist.EventID,
		wishlist.UserID,
		wishlist.Text,
	)

	return err
}

func (r *WishlistRepository) GetByUser(ctx context.Context, eventID, userID string) (*domain.Wishlist, error) {

	query := `
	SELECT id,event_id,user_id,text,created_at
	FROM wishlists
	WHERE event_id=$1 AND user_id=$2
	`

	row := r.db.QueryRow(ctx, query, eventID, userID)

	var w domain.Wishlist

	err := row.Scan(
		&w.ID,
		&w.EventID,
		&w.UserID,
		&w.Text,
		&w.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *WishlistRepository) Update(ctx context.Context, id string, text *string) error {

	query := "UPDATE wishlists SET "
	args := []interface{}{}
	argID := 1

	if text != nil {
		query += "text = $" + strconv.Itoa(argID)
		args = append(args, *text)
		argID++
	}

	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}
