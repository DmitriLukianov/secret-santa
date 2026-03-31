package assignment

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createAssignmentQuery() squirrel.InsertBuilder {
	return qb.Insert("assignments").
		Columns("id", "event_id", "giver_id", "receiver_id", "created_at")
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
