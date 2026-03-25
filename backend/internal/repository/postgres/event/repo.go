package event

import (
	"context"
	"strconv"
	"strings"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create
func (r *Repository) Create(ctx context.Context, e entity.Event) error {
	query := `
		INSERT INTO events (
			id, title, description, rules, recommendations, organizer_id,
			start_date, draw_date, end_date, status, max_participants,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(ctx, query,
		e.ID,
		e.Title,
		e.Description,
		e.Rules,
		e.Recommendations,
		e.OrganizerID,
		e.StartDate,
		e.DrawDate,
		e.EndDate,
		e.Status,
		e.MaxParticipants,
		e.CreatedAt,
		e.UpdatedAt,
	)
	return err
}

// GetByID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	query := `
		SELECT id, title, description, rules, recommendations, organizer_id,
		       start_date, draw_date, end_date, status, max_participants,
		       created_at, updated_at
		FROM events WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var e entity.Event
	err := row.Scan(
		&e.ID, &e.Title, &e.Description, &e.Rules, &e.Recommendations,
		&e.OrganizerID, &e.StartDate, &e.DrawDate, &e.EndDate,
		&e.Status, &e.MaxParticipants,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetAll
func (r *Repository) GetAll(ctx context.Context) ([]entity.Event, error) {
	query := `
		SELECT id, title, description, rules, recommendations, organizer_id,
		       start_date, draw_date, end_date, status, max_participants,
		       created_at, updated_at
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
			&e.ID, &e.Title, &e.Description, &e.Rules, &e.Recommendations,
			&e.OrganizerID, &e.StartDate, &e.DrawDate, &e.EndDate,
			&e.Status, &e.MaxParticipants,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

// Update — partial update (полная реализация)
func (r *Repository) Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error {
	query := "UPDATE events SET updated_at = NOW(), "
	args := []interface{}{}
	argID := 1

	if input.Title != nil {
		query += "title = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.Title)
		argID++
	}
	if input.Description != nil {
		query += "description = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.Description)
		argID++
	}
	if input.Rules != nil {
		query += "rules = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.Rules)
		argID++
	}
	if input.Recommendations != nil {
		query += "recommendations = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.Recommendations)
		argID++
	}
	if input.StartDate != nil {
		query += "start_date = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.StartDate)
		argID++
	}
	if input.DrawDate != nil {
		query += "draw_date = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.DrawDate)
		argID++
	}
	if input.EndDate != nil {
		query += "end_date = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.EndDate)
		argID++
	}
	if input.Status != nil {
		query += "status = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.Status)
		argID++
	}
	if input.MaxParticipants != nil {
		query += "max_participants = $" + strconv.Itoa(argID) + ", "
		args = append(args, *input.MaxParticipants)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

// Delete
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}
