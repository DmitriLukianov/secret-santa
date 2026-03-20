package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"secret-santa-backend/internal/domain"
)

type AssignmentRepository struct {
	db *pgxpool.Pool
}

func NewAssignmentRepository(db *pgxpool.Pool) *AssignmentRepository {
	return &AssignmentRepository{db: db}
}

func (r *AssignmentRepository) Create(ctx context.Context, a domain.Assignment) error {
	query := `
	INSERT INTO assignments (event_id, giver_id, receiver_id)
	VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(ctx, query, a.EventID, a.GiverID, a.ReceiverID)
	return err
}

func (r *AssignmentRepository) GetByGiver(ctx context.Context, giverID string) (*domain.Assignment, error) {
	query := `
	SELECT id, event_id, giver_id, receiver_id, created_at
	FROM assignments
	WHERE giver_id = $1
	`

	row := r.db.QueryRow(ctx, query, giverID)

	var a domain.Assignment
	err := row.Scan(&a.ID, &a.EventID, &a.GiverID, &a.ReceiverID, &a.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
func (r *AssignmentRepository) GetByEvent(ctx context.Context, eventID string) ([]domain.Assignment, error) {
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

	var result []domain.Assignment

	for rows.Next() {
		var a domain.Assignment
		err := rows.Scan(&a.ID, &a.EventID, &a.GiverID, &a.ReceiverID, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, nil
}
