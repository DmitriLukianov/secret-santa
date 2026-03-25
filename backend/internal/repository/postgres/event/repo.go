package event

import (
	"context"
	"time"

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
	query := createEventQuery().
		Values(
			event.ID, event.Name, event.Description, event.Rules, event.Recommendations,
			event.OrganizerID, event.StartDate, event.DrawDate, event.EndDate,
			event.Status, event.MaxParticipants,
		).
		Suffix("RETURNING *")

	row := r.db.QueryRowxContext(ctx, query)
	return mapRowToEvent(row)
}

func (r *Repository) GetByID(ctx context.Context, id string) (entity.Event, error) {
	uid, _ := uuid.Parse(id)
	query := getEventQuery().Where(squirrel.Eq{"id": uid})

	row := r.db.QueryRowxContext(ctx, query)
	return mapRowToEvent(row)
}

func (r *Repository) List(ctx context.Context) ([]entity.Event, error) {
	query := getEventQuery().OrderBy("created_at DESC")

	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return mapRowsToEvents(rows)
}

func (r *Repository) Update(ctx context.Context, event entity.Event) (entity.Event, error) {
	query := updateEventQuery().
		Set("name", event.Name).
		Set("description", event.Description).
		Set("rules", event.Rules).
		Set("recommendations", event.Recommendations).
		Set("start_date", event.StartDate).
		Set("draw_date", event.DrawDate).
		Set("end_date", event.EndDate).
		Set("status", event.Status).
		Set("max_participants", event.MaxParticipants).
		Set("updated_at", time.Now()).
		Where(squirrel.Eq{"id": event.ID}).
		Suffix("RETURNING *")

	row := r.db.QueryRowxContext(ctx, query)
	return mapRowToEvent(row)
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	uid, _ := uuid.Parse(id)
	query := squirrel.Delete("events").Where(squirrel.Eq{"id": uid})

	_, err := r.db.ExecContext(ctx, query)
	return err
}
