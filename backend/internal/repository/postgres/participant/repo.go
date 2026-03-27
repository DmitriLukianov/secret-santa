package participant

import (
	"context"
	"fmt"

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

func (r *Repository) Create(ctx context.Context, p entity.Participant) error {
	query := createParticipantQuery().
		Values(p.ID, p.EventID, p.UserID, p.Role, p.GiftSent, p.CreatedAt, p.UpdatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
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

func (r *Repository) UpdateGiftSent(ctx context.Context, id uuid.UUID, sent bool) error {
	query := qb.Update("participants").
		Set("gift_sent", sent).
		Set("updated_at", squirrel.Expr("NOW()"))
	if sent {
		query = query.Set("gift_sent_at", squirrel.Expr("NOW()"))
	}

	sql, args, err := query.Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
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
