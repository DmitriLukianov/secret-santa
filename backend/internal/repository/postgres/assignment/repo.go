package assignment

import (
	"context"

	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateMany(ctx context.Context, assignments []entity.Assignment) error {
	query := `
		INSERT INTO assignments (id, event_id, giver_id, receiver_id)
		VALUES ($1, $2, $3, $4)
	`

	for _, a := range assignments {
		_, err := r.db.Exec(
			ctx,
			query,
			a.ID,
			a.EventID,
			a.GiverID,
			a.ReceiverID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetByEvent(ctx context.Context, eventID string) ([]entity.Assignment, error) {
	query := `
		SELECT id, event_id, giver_id, receiver_id, created_at
		FROM assignments
		WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entity.Assignment

	for rows.Next() {
		var a entity.Assignment

		err := rows.Scan(
			&a.ID,
			&a.EventID,
			&a.GiverID,
			&a.ReceiverID,
			&a.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, a)
	}

	return result, nil
}
