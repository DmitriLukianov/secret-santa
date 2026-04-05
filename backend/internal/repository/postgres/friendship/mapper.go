package friendship

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func scanFriendship(row pgx.Row) (*entity.Friendship, error) {
	var f entity.Friendship
	err := row.Scan(
		&f.ID,
		&f.RequesterID,
		&f.AddresseeID,
		&f.Status,
		&f.CreatedAt,
		&f.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func scanFriendships(rows pgx.Rows) ([]entity.Friendship, error) {
	var result []entity.Friendship
	for rows.Next() {
		var f entity.Friendship
		err := rows.Scan(
			&f.ID,
			&f.RequesterID,
			&f.AddresseeID,
			&f.Status,
			&f.CreatedAt,
			&f.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, f)
	}
	return result, nil
}
