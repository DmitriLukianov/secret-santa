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
	return qb.Select("id", "event_id", "giver_id", "receiver_id", "created_at").
		From("assignments").
		Where(squirrel.Eq{"event_id": eventID})
}

func deleteAssignmentsByEventQuery(eventID string) squirrel.DeleteBuilder {
	return qb.Delete("assignments").
		Where(squirrel.Eq{"event_id": eventID})
}
