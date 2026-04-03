package wishlist

import (
	"context"
	"fmt"
	"time"

	"secret-santa-backend/internal/definitions"
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

	existing, err := r.GetByParticipant(ctx, w.ParticipantID)
	if err == nil && existing != nil {
		return fmt.Errorf("wishlist already exists for participant %s: %w", w.ParticipantID, definitions.ErrConflict)
	}

	query := createWishlistQuery().
		Values(w.ID, w.ParticipantID, w.Visibility, w.CreatedAt, w.UpdatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) CreateItem(ctx context.Context, item entity.WishlistItem) error {
	query := createWishlistItemQuery().
		Values(item.ID, item.WishlistID, item.Title, item.Link, item.ImageURL, item.Comment, item.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetByParticipant(ctx context.Context, participantID uuid.UUID) (*entity.Wishlist, error) {
	query := getWishlistByParticipantQuery(participantID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanWishlist(row)
}

func (r *Repository) GetItems(ctx context.Context, wishlistID uuid.UUID) ([]entity.WishlistItem, error) {
	query := getWishlistItemsQuery(wishlistID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanWishlistItems(rows)
}
func (r *Repository) UpdateItem(ctx context.Context, itemID uuid.UUID, title string, link, imageURL, comment *string) error {
	query := updateWishlistItemQuery(itemID.String()).
		Set("title", title).
		Set("link", link).
		Set("image_url", imageURL).
		Set("comment", comment).
		Set("created_at", time.Now())

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	query := deleteWishlistItemQuery(itemID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetItemByID(ctx context.Context, itemID uuid.UUID) (*entity.WishlistItem, error) {
	if itemID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	query := getWishlistItemByIDQuery(itemID)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	return scanWishlistItem(row)
}

func (r *Repository) GetByID(ctx context.Context, wishlistID uuid.UUID) (*entity.Wishlist, error) {
	if wishlistID == uuid.Nil {
		return nil, definitions.ErrInvalidUserInput
	}

	query := getWishlistByIDQuery(wishlistID)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	return scanWishlist(row)
}
