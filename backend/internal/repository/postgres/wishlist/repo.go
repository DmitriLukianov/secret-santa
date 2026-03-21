package wishlist

import (
	"context"

	"secret-santa-backend/internal/entity"

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
		INSERT INTO wishlists (
			id, user_id, title, description, link, image_url, visibility
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		w.ID,
		w.UserID,
		w.Title,
		w.Description,
		w.Link,
		w.ImageURL,
		w.Visibility,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*entity.Wishlist, error) {
	query := `
		SELECT id, user_id, title, description, link, image_url, visibility, created_at
		FROM wishlists
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var w entity.Wishlist

	err := row.Scan(
		&w.ID,
		&w.UserID,
		&w.Title,
		&w.Description,
		&w.Link,
		&w.ImageURL,
		&w.Visibility,
		&w.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &w, nil
}

func (r *Repository) GetByUser(ctx context.Context, userID string) ([]entity.Wishlist, error) {
	query := `
		SELECT id, user_id, title, description, link, image_url, visibility, created_at
		FROM wishlists
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Wishlist

	for rows.Next() {
		var w entity.Wishlist

		if err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.Title,
			&w.Description,
			&w.Link,
			&w.ImageURL,
			&w.Visibility,
			&w.CreatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, w)
	}

	return result, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	title, description, link, imageURL, visibility *string,
) error {
	query := `
		UPDATE wishlists
		SET
			title = COALESCE($1, title),
			description = COALESCE($2, description),
			link = COALESCE($3, link),
			image_url = COALESCE($4, image_url),
			visibility = COALESCE($5, visibility)
		WHERE id = $6
	`

	_, err := r.db.Exec(
		ctx,
		query,
		title,
		description,
		link,
		imageURL,
		visibility,
		id,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM wishlists WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
