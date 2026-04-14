package invitation

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func ScanInvitation(row pgx.Row) (*entity.Invitation, error) {
	var i entity.Invitation
	err := row.Scan(
		&i.ID,
		&i.EventID,
		&i.Token,
		&i.ExpiresAt,
		&i.CreatedBy,
		&i.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}
