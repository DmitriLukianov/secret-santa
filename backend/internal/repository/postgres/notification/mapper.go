package notification

import (
	"encoding/json"

	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanNotification(row pgx.Row) (*entity.Notification, error) {
	var n entity.Notification
	var payloadRaw []byte
	err := row.Scan(
		&n.ID,
		&n.UserID,
		&n.Type,
		&payloadRaw,
		&n.IsRead,
		&n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(payloadRaw, &n.Payload); err != nil {
		n.Payload = map[string]string{}
	}
	return &n, nil
}

func scanNotifications(rows pgx.Rows) ([]entity.Notification, error) {
	var result []entity.Notification
	for rows.Next() {
		var n entity.Notification
		var payloadRaw []byte
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&payloadRaw,
			&n.IsRead,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payloadRaw, &n.Payload); err != nil {
			n.Payload = map[string]string{}
		}
		result = append(result, n)
	}
	return result, nil
}
