package postgres

import (
	"context"
	"secret-santa-backend/internal/domain"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) CreateEvent(ctx context.Context, event domain.Event) error {

	query := `
	INSERT INTO events (name, description, organizer_id, start_date, draw_date, end_date)
	VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		event.Name,
		event.Description,
		event.OrganizerID,
		event.StartDate,
		event.DrawDate,
		event.EndDate,
	)

	return err
}

func (r *EventRepository) GetEventByID(ctx context.Context, id string) (*domain.Event, error) {

	query := `
	SELECT id, name, description, organizer_id, start_date, draw_date, end_date, created_at
	FROM events
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var event domain.Event

	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.OrganizerID,
		&event.StartDate,
		&event.DrawDate,
		&event.EndDate,
		&event.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &event, nil
}
func (r *EventRepository) GetEvents(ctx context.Context) ([]domain.Event, error) {

	query := `
	SELECT id, name, description, organizer_id, start_date, draw_date, end_date, created_at
	FROM events
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []domain.Event

	for rows.Next() {

		var event domain.Event

		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description,
			&event.OrganizerID,
			&event.StartDate,
			&event.DrawDate,
			&event.EndDate,
			&event.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) UpdateEvent(ctx context.Context, id string, name, description *string) error {

	query := "UPDATE events SET "
	args := []interface{}{}
	argID := 1

	if name != nil {
		query += "name = $" + strconv.Itoa(argID) + ", "
		args = append(args, *name)
		argID++
	}

	if description != nil {
		query += "description = $" + strconv.Itoa(argID) + ", "
		args = append(args, *description)
		argID++
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *EventRepository) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
