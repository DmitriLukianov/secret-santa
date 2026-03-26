package wishlist

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, w entity.Wishlist) error {
	query := `
		INSERT INTO wishlists (id, participant_id, visibility, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		w.ID, w.ParticipantID, w.Visibility, w.CreatedAt, w.UpdatedAt,
	)
	return err
}

func (r *Repository) CreateItem(ctx context.Context, item entity.WishlistItem) error {
	query := `
		INSERT INTO wishlist_items (id, wishlist_id, title, link, image_url, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		item.ID, item.WishlistID, item.Title, item.Link, item.ImageURL, item.Comment, item.CreatedAt,
	)
	return err
}

func (r *Repository) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	query := `
		SELECT id, participant_id, visibility, created_at, updated_at
		FROM wishlists WHERE participant_id = $1
	`

	row := r.db.QueryRow(ctx, query, participantID)
	return ScanWishlist(row)
}

func (r *Repository) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	query := `
		SELECT id, wishlist_id, title, link, image_url, comment, created_at
		FROM wishlist_items WHERE wishlist_id = $1
	`

	rows, err := r.db.Query(ctx, query, wishlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanWishlistItems(rows)
}
