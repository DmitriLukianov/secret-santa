package event

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, event entity.Event) (entity.Event, error) {
	query, args, err := createEventQuery().
		Values(
			event.ID, event.Name, event.Description, event.Rules, event.Recommendations,
			event.OrganizerID, event.StartDate, event.DrawDate, event.EndDate,
			event.Status, event.MaxParticipants,
		).
		Suffix("RETURNING id, name, description, rules, recommendations, organizer_id, start_date, draw_date, end_date, status, max_participants, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	return mapRowToEvent(row)
}

func (r *Repository) GetByID(ctx context.Context, id string) (entity.Event, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return entity.Event{}, err
	}

	query, args, err := getEventQuery().Where(squirrel.Eq{"id": uid}).ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	return mapRowToEvent(row)
}

func (r *Repository) List(ctx context.Context) ([]entity.Event, error) {
	query, args, err := listEventsQuery().OrderBy("created_at DESC").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRowsToEvents(rows)
}

func (r *Repository) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	query, args, err := updateEventQuery().
		Set("name", event.Name).
		Set("description", event.Description).
		Set("rules", event.Rules).
		Set("recommendations", event.Recommendations).
		Set("start_date", event.StartDate).
		Set("draw_date", event.DrawDate).
		Set("end_date", event.EndDate).
		Set("status", event.Status).
		Set("max_participants", event.MaxParticipants).
		Set("updated_at", event.UpdatedAt).
		Where(squirrel.Eq{"id": event.ID}).
		Suffix("RETURNING id, name, description, rules, recommendations, organizer_id, start_date, draw_date, end_date, status, max_participants, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Event{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	return mapRowToEvent(row)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	query, args, err := deleteEventQuery().Where(squirrel.Eq{"id": uid}).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}
