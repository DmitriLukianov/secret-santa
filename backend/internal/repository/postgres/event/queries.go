package event

import "github.com/Masterminds/squirrel"

// getEventQuery — базовый запрос для получения события
func getEventQuery() squirrel.SelectBuilder {
	return squirrel.Select(
		"id", "name", "description", "rules", "recommendations",
		"organizer_id", "start_date", "draw_date", "end_date",
		"status", "max_participants", "created_at", "updated_at",
	).
		From("events")
}

// createEventQuery — запрос на создание события
func createEventQuery() squirrel.InsertBuilder {
	return squirrel.Insert("events").
		Columns(
			"id", "name", "description", "rules", "recommendations",
			"organizer_id", "start_date", "draw_date", "end_date",
			"status", "max_participants",
		)
}

// updateEventQuery — запрос на обновление события
func updateEventQuery() squirrel.UpdateBuilder {
	return squirrel.Update("events")
}
