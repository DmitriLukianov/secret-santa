package event

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func ScanEvent(row pgx.Row) (*entity.Event, error) {
	var e entity.Event
	err := row.Scan(
		&e.ID,
		&e.Title,
		&e.OrganizerNotes,
		&e.OrganizerID,
		&e.StartDate,
		&e.DrawDate,
		&e.Status,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func ScanEvents(rows pgx.Rows) ([]entity.Event, error) {
	var events []entity.Event
	for rows.Next() {
		e, err := ScanEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, *e)
	}
	return events, nil
}
