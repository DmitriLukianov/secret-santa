package chat

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanMessage(row pgx.Row) (*entity.Message, error) {
	var m entity.Message
	err := row.Scan(
		&m.ID,
		&m.EventID,
		&m.SenderID,
		&m.ReceiverID,
		&m.Content,
		&m.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func scanMessages(rows pgx.Rows) ([]entity.Message, error) {
	var messages []entity.Message
	for rows.Next() {
		var m entity.Message
		err := rows.Scan(
			&m.ID,
			&m.EventID,
			&m.SenderID,
			&m.ReceiverID,
			&m.Content,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
