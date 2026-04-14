package participant

import (
	"context"
	"fmt"

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

func (r *Repository) Create(ctx context.Context, p entity.Participant) (entity.Participant, error) {
	query, args, err := createParticipantQuery().
		Values(p.EventID, p.UserID).
		Suffix("RETURNING id, event_id, user_id, created_at").
		ToSql()

	if err != nil {
		return entity.Participant{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	returned, err := scanParticipant(row)
	if err != nil {
		return entity.Participant{}, err
	}

	return *returned, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	query := getParticipantByIDQuery(id.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanParticipant(row)
}

func (r *Repository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error) {
	query := getParticipantsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanParticipants(rows)
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := deleteParticipantQuery(id.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error) {
	query := getParticipantByUserAndEventQuery(userID.String(), eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	p, err := scanParticipant(row)
	if err != nil {
		return nil, fmt.Errorf("get participant by user and event: %w", err)
	}
	return p, nil
}
