package assignment

import (
	"context"

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

func (r *Repository) Create(ctx context.Context, a entity.Assignment) error {
	query := createAssignmentQuery().
		Values(a.ID, a.EventID, a.GiverID, a.ReceiverID, a.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error) {
	query := getAssignmentsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAssignments(rows)
}

func (r *Repository) DeleteByEvent(ctx context.Context, eventID uuid.UUID) error {
	query := deleteAssignmentsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
