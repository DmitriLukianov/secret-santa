package event

import (
	"context"

	"secret-santa-backend/internal/dto"
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

// Create
func (r *Repository) Create(ctx context.Context, e entity.Event) error {
	query, args, err := createEventQuery().
		Values(
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
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

// GetByID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Event, error) {
	query, args, err := getEventQuery().
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	return ScanEvent(row)
}

// GetAll
func (r *Repository) GetAll(ctx context.Context) ([]entity.Event, error) {
	query, args, err := listEventsQuery().ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanEvents(rows)
}

// Update
func (r *Repository) Update(ctx context.Context, id uuid.UUID, input dto.UpdateEventInput) error {
	q := updateEventQuery().Set("updated_at", "NOW()")

	if input.Title != nil {
		q = q.Set("title", *input.Title)
	}
	if input.Description != nil {
		q = q.Set("description", *input.Description)
	}
	if input.Rules != nil {
		q = q.Set("rules", *input.Rules)
	}
	if input.Recommendations != nil {
		q = q.Set("recommendations", *input.Recommendations)
	}
	if input.StartDate != nil {
		q = q.Set("start_date", *input.StartDate)
	}
	if input.DrawDate != nil {
		q = q.Set("draw_date", *input.DrawDate)
	}
	if input.EndDate != nil {
		q = q.Set("end_date", *input.EndDate)
	}
	if input.MaxParticipants != nil {
		q = q.Set("max_participants", *input.MaxParticipants)
	}

	query, args, err := q.Where("id = ?", id).ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

// UpdateStatus
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status entity.EventStatus) error {
	query, args, err := updateEventQuery().
		Set("status", status).
		Set("updated_at", "NOW()").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

// Delete
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := deleteEventQuery().
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

// GetEventsForUser — теперь тоже через queries.go
func (r *Repository) GetEventsForUser(ctx context.Context, userID uuid.UUID) ([]entity.Event, error) {
	query, args, err := getEventsForUserQuery().
		Where(squirrel.Or{
			squirrel.Eq{"events.organizer_id": userID},
			squirrel.Eq{"p.user_id": userID},
		}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanEvents(rows)
}
