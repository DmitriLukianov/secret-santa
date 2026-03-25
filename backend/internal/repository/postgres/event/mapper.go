package event

import (
	"secret-santa-backend/internal/entity"

	"github.com/georgysavva/scany/v2/dbscan"
)

// mapRowToEvent — преобразует одну строку из БД в entity.Event
func mapRowToEvent(row dbscan.Row) (entity.Event, error) {
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
	if err != nil {
		return entity.Event{}, err
	}
	return e, nil
}

// mapRowsToEvents — преобразует несколько строк в слайс событий
func mapRowsToEvents(rows dbscan.Rows) ([]entity.Event, error) {
	var events []entity.Event

	for rows.Next() {
		event, err := mapRowToEvent(rows) // ← здесь передаём Rows, но внутри функции используется Row
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
