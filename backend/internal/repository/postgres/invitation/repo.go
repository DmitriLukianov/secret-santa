package invitation

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

func (r *Repository) Create(ctx context.Context, i entity.Invitation) error {
	query, args, err := createInvitationQuery().
		Values(
			i.ID,
			i.EventID,
			i.Token,
			i.ExpiresAt,
			i.CreatedBy,
			i.CreatedAt,
			i.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) GetByToken(ctx context.Context, token string) (*entity.Invitation, error) {
	query, args, err := getInvitationByTokenQuery().
		Where("token = ?", token).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	return ScanInvitation(row)
}
