package postgres

import (
	"context"
	"secret-santa-backend/internal/domain"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) error {
	query := `
	INSERT INTO users (name, email)
	VALUES ($1, $2)
	`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email)
	return err
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
	SELECT id, name, email, created_at
	FROM users
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	var user domain.User

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
func (r *userRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User

		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id string, name, email *string) error {

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

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = $" + strconv.Itoa(argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)
	return err
}
