package user

import (
	"context"
	"errors"

	"secret-santa-backend/internal/definitions"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create создаёт пользователя и возвращает полностью заполненную сущность из БД
func (r *Repository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	query := Create().
		Values(user.Name, user.Email, user.OAuthID, user.OAuthProvider)

	sql, args, err := query.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	returnedUser, err := scanUser(row)
	if err != nil {
		return entity.User{}, err
	}

	return *returnedUser, nil
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
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return definitions.ErrEmailTaken
		}
	}
	return err
}

