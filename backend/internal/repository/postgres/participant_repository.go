package postgres

import (
	"context"
	"secret-santa-backend/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ParticipantRepository struct {
	db *pgxpool.Pool
}

func NewParticipantRepository(db *pgxpool.Pool) *ParticipantRepository {
	return &ParticipantRepository{db: db}
}

func (r *ParticipantRepository) AddParticipant(ctx context.Context, p domain.Participant) error {

	query := `
	INSERT INTO participants (event_id, user_id)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, p.EventID, p.UserID)
	return err
}

func (r *ParticipantRepository) GetParticipantsByEvent(ctx context.Context, eventID string) ([]domain.Participant, error) {

	query := `
	SELECT id, event_id, user_id, joined_at
	FROM participants
	WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.Participant

	for rows.Next() {
		var p domain.Participant

		err := rows.Scan(
			&p.ID,
			&p.EventID,
			&p.UserID,
			&p.JoinedAt,
		)
		if err != nil {
			return nil, err
		}

		participants = append(participants, p)
	}

	return participants, nil
}

func (r *ParticipantRepository) DeleteParticipant(ctx context.Context, eventID, userID string) error {

	query := `
	DELETE FROM participants
	WHERE event_id = $1 AND user_id = $2
	`

	_, err := r.db.Exec(ctx, query, eventID, userID)
	return err
}
