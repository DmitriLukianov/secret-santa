package user

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

func (r *Repository) Create(ctx context.Context, user entity.User) error {
	query := Create().
		Values(user.Name, user.Email, user.OAuthID, user.OAuthProvider)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := GetByID(id)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanUser(row)
}

func (r *Repository) GetByOAuthID(ctx context.Context, oauthID, oauthProvider string) (*entity.User, error) {
	query := GetByOAuthID(oauthID, oauthProvider)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanUser(row)
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := GetByEmail(email)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanUser(row)
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.User, error) {
	query := GetAll()
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, name, email *string) error {
	query := Update(id)
	if name != nil {
		query = query.Set("name", *name)
	}
	if email != nil {
		query = query.Set("email", *email)
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := Delete(id)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
