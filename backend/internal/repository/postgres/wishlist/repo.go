package wishlist

import (
	"context"
	"fmt"

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

func (r *Repository) Create(ctx context.Context, w entity.Wishlist) (entity.Wishlist, error) {
	// Проверка на существование (можно оставить, но лучше делать через unique constraint)
	existing, err := r.GetByParticipant(ctx, w.ParticipantID)
	if err == nil && existing != nil {
		return entity.Wishlist{}, fmt.Errorf("wishlist already exists for participant %s: %w", w.ParticipantID, definitions.ErrConflict)
	}

	query, args, err := createWishlistQuery().
		Values(w.ParticipantID, w.Visibility).
		Suffix("RETURNING id, participant_id, visibility, created_at, updated_at").
		ToSql()

	if err != nil {
		return entity.Wishlist{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	returned, err := scanWishlist(row)
	if err != nil {
		return entity.Wishlist{}, err
	}

	return *returned, nil
}

func (r *Repository) CreateItem(ctx context.Context, item entity.WishlistItem) (entity.WishlistItem, error) {
	query, args, err := createWishlistItemQuery().
		Values(item.WishlistID, item.Title, item.Link, item.ImageURL, item.Comment).
		Suffix("RETURNING id, wishlist_id, title, link, image_url, comment, created_at").
		ToSql()

	if err != nil {
		return entity.WishlistItem{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	returned, err := scanWishlistItem(row)
	if err != nil {
		return entity.WishlistItem{}, err
	}

	return *returned, nil
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
		Set("updated_at", "NOW()") // добавили обновление updated_at

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
