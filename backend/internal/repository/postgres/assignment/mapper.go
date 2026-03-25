package assignment

import (
	"secret-santa-backend/internal/entity"

	"github.com/jackc/pgx/v5"
)

func ScanAssignment(row pgx.Row) (*entity.Assignment, error) {
	var a entity.Assignment
	err := row.Scan(
		&a.ID,
		&a.EventID,
		&a.GiverID,
		&a.ReceiverID,
		&a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func ScanAssignments(rows pgx.Rows) ([]entity.Assignment, error) {
	var assignments []entity.Assignment
	for rows.Next() {
		a, err := ScanAssignment(rows)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, *a)
	}
	return assignments, nil
}
