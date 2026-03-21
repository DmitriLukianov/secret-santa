package user

import (
	"context"
	"strconv"
	"strings"

	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CREATE
func (r *Repository) Create(ctx context.Context, user entity.User) error {
	query := `
		INSERT INTO users (id, name, email)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(ctx, query,
		user.ID,
		user.Name,
		user.Email,
	)

	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]entity.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User

	for rows.Next() {
		var u entity.User

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *Repository) Update(ctx context.Context, id string, name, email *string) error {
	query := "UPDATE users SET "
	args := []interface{}{}
	argID := 1

	if name != nil {
		query += "name = $" + strconv.Itoa(argID) + ", "
		args = append(args, *name)
		argID++
	}

	if email != nil {
		query += "email = $" + strconv.Itoa(argID) + ", "
		args = append(args, *email)
		argID++
	}

	if len(args) == 0 {
		return nil
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
