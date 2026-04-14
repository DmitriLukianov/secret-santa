package participant

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanParticipant(row pgx.Row) (*entity.Participant, error) {
	var p entity.Participant
	err := row.Scan(
		&p.ID,
		&p.EventID,
		&p.UserID,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func scanParticipants(rows pgx.Rows) ([]entity.Participant, error) {
	var participants []entity.Participant
	for rows.Next() {
		var p entity.Participant
		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.UserID,
			&p.CreatedAt,
			&p.UserName,
			&p.UserEmail,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}
