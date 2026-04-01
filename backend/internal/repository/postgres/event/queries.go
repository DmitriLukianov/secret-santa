package event

import "github.com/Masterminds/squirrel"

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func getEventQuery() squirrel.SelectBuilder {
	return psql.Select(
		"id", "title", "description", "rules", "recommendations",
		"organizer_id", "start_date", "draw_date", "end_date",
		"status", "max_participants", "created_at", "updated_at",
	).
		From("events")
}

func listEventsQuery() squirrel.SelectBuilder {
	return getEventQuery()
}

func createEventQuery() squirrel.InsertBuilder {
	return psql.Insert("events").
		Columns(
			"id", "title", "description", "rules", "recommendations",
			"organizer_id", "start_date", "draw_date", "end_date",
			"status", "max_participants",
		)
}

func updateEventQuery() squirrel.UpdateBuilder {
	return psql.Update("events")
}

func deleteEventQuery() squirrel.DeleteBuilder {
	return psql.Delete("events")
}

func getEventsForUserQuery() squirrel.SelectBuilder {
	return getEventQuery().
		Distinct().
		LeftJoin("participants p ON events.id = p.event_id").
		Where(squirrel.Or{
			squirrel.Eq{"events.organizer_id": "?"},
			squirrel.Eq{"p.user_id": "?"},
		}).
		OrderBy("events.created_at DESC")
}
