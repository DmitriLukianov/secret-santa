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

func (r *Repository) Create(ctx context.Context, p entity.Participant) error {
	query := `
		INSERT INTO participants (id, event_id, user_id, role, gift_sent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		p.ID, p.EventID, p.UserID, p.Role, p.GiftSent,
		p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Participant, error) {
	query := `
		SELECT id, event_id, user_id, role, gift_sent, gift_sent_at, created_at, updated_at
		FROM participants WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	return ScanParticipant(row)
}

func (r *Repository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Participant, error) {
	query := `
		SELECT id, event_id, user_id, role, gift_sent, gift_sent_at, created_at, updated_at
		FROM participants WHERE event_id = $1
	`

	rows, err := r.db.Query(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return ScanParticipants(rows)
}

func (r *Repository) UpdateGiftSent(ctx context.Context, id uuid.UUID, sent bool) error {
	query := `
		UPDATE participants 
		SET gift_sent = $1, gift_sent_at = NOW(), updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, sent, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM participants WHERE id = $1`, id)
	return err
}

func (r *Repository) GetByUserAndEvent(ctx context.Context, userID, eventID uuid.UUID) (*entity.Participant, error) {
	query := `
		SELECT id, event_id, user_id, role, gift_sent, gift_sent_at, created_at
		FROM participants 
		WHERE user_id = $1 AND event_id = $2
		LIMIT 1
	`

	var p entity.Participant
	err := r.db.QueryRow(ctx, query, userID, eventID).
		Scan(&p.ID, &p.EventID, &p.UserID, &p.Role, &p.GiftSent, &p.GiftSentAt, &p.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("get participant by user and event: %w", err)
	}

	return &p, nil
}
