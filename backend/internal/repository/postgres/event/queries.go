package event

import "github.com/Masterminds/squirrel"

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// getEventQuery — ВСЕ колонки теперь с префиксом events.
func getEventQuery() squirrel.SelectBuilder {
	return psql.Select(
		"events.id",
		"events.title",
		"events.description",
		"events.rules",
		"events.recommendations",
		"events.organizer_id",
		"events.start_date",
		"events.draw_date",
		"events.end_date",
		"events.status",
		"events.max_participants",
		"events.created_at",
		"events.updated_at",
	).
		From("events")
}

func listEventsQuery() squirrel.SelectBuilder {
	return getEventQuery()
}

// createEventQuery — без изменений (БД сама даёт id, created_at, updated_at)
func createEventQuery() squirrel.InsertBuilder {
	return psql.Insert("events").
		Columns(
			"title", "description", "rules", "recommendations",
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

// getEventsForUserQuery — JOIN остался, но теперь SELECT корректный
func getEventsForUserQuery() squirrel.SelectBuilder {
	return getEventQuery().
		Distinct().
		LeftJoin("participants p ON events.id = p.event_id").
		OrderBy("events.created_at DESC")
}
