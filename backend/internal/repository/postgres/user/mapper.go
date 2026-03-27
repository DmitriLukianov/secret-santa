package user

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanUser(row pgx.Row) (*entity.User, error) {
	var u entity.User
	err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.OAuthID,
		&u.OAuthProvider,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func scanUsers(rows pgx.Rows) ([]entity.User, error) {
	var users []entity.User
	for rows.Next() {
		var u entity.User
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.OAuthID,
			&u.OAuthProvider,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
