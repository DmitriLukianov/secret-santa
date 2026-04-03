package verification

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type VerificationCode struct {
	ID        uuid.UUID
	Email     string
	Code      string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

func scanVerificationCode(row pgx.Row) (*VerificationCode, error) {
	var vc VerificationCode
	err := row.Scan(
		&vc.ID,
		&vc.Email,
		&vc.Code,
		&vc.ExpiresAt,
		&vc.Used,
		&vc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &vc, nil
}
