package friendship

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

func (r *Repository) Create(ctx context.Context, f entity.Friendship) (entity.Friendship, error) {
	query, args, err := createFriendshipQuery().
		Values(f.RequesterID, f.AddresseeID, f.Status).
		Suffix("RETURNING id, requester_id, addressee_id, status, created_at, updated_at").
		ToSql()
	if err != nil {
		return entity.Friendship{}, err
	}
	row := r.db.QueryRow(ctx, query, args...)
	created, err := scanFriendship(row)
	if err != nil {
		return entity.Friendship{}, err
	}
	return *created, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Friendship, error) {
	sql, args, err := getFriendshipByIDQuery(id).ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanFriendship(row)
}

func (r *Repository) GetByUsers(ctx context.Context, userA, userB uuid.UUID) (*entity.Friendship, error) {
	sql, args, err := getFriendshipByUsersQuery(userA, userB).ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.QueryRow(ctx, sql, args...)
	return scanFriendship(row)
}

func (r *Repository) GetFriends(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error) {
	sql, args, err := getFriendsQuery(userID).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFriendships(rows)
}

func (r *Repository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]entity.Friendship, error) {
	sql, args, err := getPendingRequestsQuery(userID).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanFriendships(rows)
}

func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	sql, args, err := updateFriendshipStatusQuery(id, status).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := deleteFriendshipQuery(id).ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}
