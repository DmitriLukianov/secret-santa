package chat

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

func (r *Repository) CreateMessage(ctx context.Context, msg entity.Message) error {
	query := createMessageQuery().
		Values(msg.ID, msg.EventID, msg.SenderID, msg.ReceiverID, msg.Content, msg.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetMessagesByPair(ctx context.Context, eventID, user1ID, user2ID uuid.UUID) ([]entity.Message, error) {
	query := getMessagesByPairQuery(eventID.String(), user1ID.String(), user2ID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanMessages(rows)
}
