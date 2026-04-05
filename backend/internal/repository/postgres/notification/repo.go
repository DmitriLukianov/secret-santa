package notification

import (
	"context"
	"encoding/json"

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

func (r *Repository) Create(ctx context.Context, n entity.Notification) (entity.Notification, error) {
	payloadJSON, err := json.Marshal(n.Payload)
	if err != nil {
		return entity.Notification{}, err
	}

	query, args, err := createNotificationQuery().
		Values(n.UserID, n.Type, payloadJSON).
		Suffix("RETURNING id, user_id, type, payload, is_read, created_at").
		ToSql()
	if err != nil {
		return entity.Notification{}, err
	}

	row := r.db.QueryRow(ctx, query, args...)
	created, err := scanNotification(row)
	if err != nil {
		return entity.Notification{}, err
	}
	return *created, nil
}

func (r *Repository) GetByUser(ctx context.Context, userID uuid.UUID) ([]entity.Notification, error) {
	sql, args, err := getNotificationsByUserQuery(userID).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanNotifications(rows)
}

func (r *Repository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	sql, args, err := markAsReadQuery(id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) MarkAllAsRead(ctx context.Context, userID uuid.UUID) error {
	sql, args, err := markAllAsReadQuery(userID).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
