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

// Create
func (r *Repository) Create(ctx context.Context, a entity.Assignment) error {
	query := `
		INSERT INTO assignments (id, event_id, giver_id, receiver_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		a.ID, a.EventID, a.GiverID, a.ReceiverID, a.CreatedAt,
	)
	return err
}

// GetByEvent
func (r *Repository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error) {
	query := `
		SELECT id, event_id, giver_id, receiver_id, created_at
		FROM assignments WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanAssignments(rows)
}

// DeleteByEvent (для пересоздания жеребьёвки)
func (r *Repository) DeleteByEvent(ctx context.Context, eventID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM assignments WHERE event_id = $1`, eventID)
	return err
}
