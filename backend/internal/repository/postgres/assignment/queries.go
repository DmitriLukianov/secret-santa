package assignment

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// createAssignmentQuery — теперь DB-first (убрали id и created_at)
func createAssignmentQuery() squirrel.InsertBuilder {
	return qb.Insert("assignments").
		Columns("event_id", "giver_id", "receiver_id")
	// id и created_at генерирует PostgreSQL (DEFAULT gen_random_uuid() и NOW())
}

func getAssignmentsByEventQuery(eventID string) squirrel.SelectBuilder {
	return qb.Select("a.id", "a.event_id", "a.giver_id", "a.receiver_id", "u.name AS receiver_name", "a.created_at").
		From("assignments a").
		Join("users u ON u.id = a.receiver_id").
		Where(squirrel.Eq{"a.event_id": eventID})
}

func deleteAssignmentsByEventQuery(eventID string) squirrel.DeleteBuilder {
	return qb.Delete("assignments").
		Where(squirrel.Eq{"event_id": eventID})
}
