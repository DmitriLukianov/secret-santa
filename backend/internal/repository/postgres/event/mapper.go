package event

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

// ScanEvent сканирует одну строку из БД в entity.Event
func ScanEvent(row pgx.Row) (*entity.Event, error) {
	var e entity.Event
	err := row.Scan(
		&e.ID,
		&e.Title,
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
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// ScanEvents сканирует несколько строк
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
