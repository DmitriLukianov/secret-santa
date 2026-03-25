package event

import (
	"github.com/jackc/pgx/v5"

	"secret-santa-backend/internal/entity"
)

func mapRowToEvent(row pgx.Row) (entity.Event, error) {
	var e entity.Event
	err := row.Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.Rules,
		&e.Recommendations,
		&e.OrganizerID,
		&e.StartDate,
		&e.DrawDate,
		&e.EndDate,
		&e.Status,
		&e.MaxParticipants,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	return e, err
}

func mapRowsToEvents(rows pgx.Rows) ([]entity.Event, error) {
	var events []entity.Event

	for rows.Next() {
		var e entity.Event
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.Rules,
			&e.Recommendations,
			&e.OrganizerID,
			&e.StartDate,
			&e.DrawDate,
			&e.EndDate,
			&e.Status,
			&e.MaxParticipants,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}
