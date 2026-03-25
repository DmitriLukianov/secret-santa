package participant

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func ScanParticipant(row pgx.Row) (*entity.Participant, error) {
	var p entity.Participant
	err := row.Scan(
		&p.ID,
		&p.EventID,
		&p.UserID,
		&p.Role,
		&p.GiftSent,
		&p.GiftSentAt,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func ScanParticipants(rows pgx.Rows) ([]entity.Participant, error) {
	var participants []entity.Participant
	for rows.Next() {
		p, err := ScanParticipant(rows)
		if err != nil {
			return nil, err
		}
		participants = append(participants, *p)
	}
	return participants, nil
}
