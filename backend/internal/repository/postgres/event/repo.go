package event

import (
	"context"
	"strconv"
	"strings"

	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CREATE
func (r *Repository) Create(ctx context.Context, event entity.Event) error {
	query := `
		INSERT INTO events (id, name, description, organizer_id, start_date, draw_date, end_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`

	_, err := r.db.Exec(ctx, query,
		event.ID, // 👈 ВАЖНО (добавили id)
		event.Name,
		event.Description,
		event.OrganizerID,
		event.StartDate,
		event.DrawDate,
		event.EndDate,
	)

	return err
}

// GET BY ID
func (r *Repository) GetByID(ctx context.Context, id string) (*entity.Event, error) {
	query := `
		SELECT id, name, description, organizer_id, start_date, draw_date, end_date, created_at
		FROM events
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var e entity.Event
	err := row.Scan(
		&e.ID,
		&e.Name,
		&e.Description,
		&e.OrganizerID,
		&e.StartDate,
		&e.DrawDate,
		&e.EndDate,
		&e.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

// GET ALL
func (r *Repository) GetAll(ctx context.Context) ([]entity.Event, error) {
	query := `
		SELECT id, name, description, organizer_id, start_date, draw_date, end_date, created_at
		FROM events
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []entity.Event

	for rows.Next() {
		var e entity.Event

		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Description,
			&e.OrganizerID,
			&e.StartDate,
			&e.DrawDate,
			&e.EndDate,
			&e.CreatedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

// UPDATE (partial)
func (r *Repository) Update(ctx context.Context, id string, name, description *string) error {
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

	// если ничего не передали
	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

// DELETE
func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
