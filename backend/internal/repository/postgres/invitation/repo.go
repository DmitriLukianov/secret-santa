package invitation

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, i entity.Invitation) error {
	query := `
		INSERT INTO invitations (
			id, event_id, token, expires_at, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		i.ID, i.EventID, i.Token, i.ExpiresAt,
		i.CreatedBy, i.CreatedAt, i.UpdatedAt,
	)
	return err
}

func (r *Repository) GetByToken(ctx context.Context, token string) (*entity.Invitation, error) {
	query := `
		SELECT id, event_id, token, expires_at, created_by, created_at, updated_at
		FROM invitations
		WHERE token = $1
	`

	row := r.db.QueryRow(ctx, query, token)
	return scanInvitation(row)
}

// scanInvitation — обновлённая версия без поля Used
func scanInvitation(row pgx.Row) (*entity.Invitation, error) {
	var i entity.Invitation
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan invitation: %w", err)
	}
	return &i, nil
}
