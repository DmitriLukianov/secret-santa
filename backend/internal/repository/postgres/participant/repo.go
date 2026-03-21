package participant

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

func (r *Repository) GetByEvent(ctx context.Context, eventID string) ([]entity.Participant, error) {
	query := `
		SELECT id, event_id, user_id
		FROM participants
		WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []entity.Participant

	for rows.Next() {
		var p entity.Participant

		if err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.UserID,
		); err != nil {
			return nil, err
		}

		participants = append(participants, p)
	}

	return participants, nil
}

func (r *Repository) Add(ctx context.Context, p entity.Participant) error {
	query := `
		INSERT INTO participants (id, event_id, user_id)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query,
		p.ID,
		p.EventID,
		p.UserID,
	)

	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM participants WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
