package participant

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createParticipantQuery() squirrel.InsertBuilder {
	return qb.Insert("participants").
		Columns("event_id", "user_id", "role", "gift_sent")
	// created_at, updated_at, id — пусть БД сама заполняет
}

func getParticipantByIDQuery(id string) squirrel.SelectBuilder {
	return qb.Select("id", "event_id", "user_id", "role", "gift_sent", "gift_sent_at", "created_at", "updated_at").
		From("participants").
		Where(squirrel.Eq{"id": id})
}

func getParticipantsByEventQuery(eventID string) squirrel.SelectBuilder {
	return qb.Select("id", "event_id", "user_id", "role", "gift_sent", "gift_sent_at", "created_at", "updated_at").
		From("participants").
		Where(squirrel.Eq{"event_id": eventID})
}

func deleteParticipantQuery(id string) squirrel.DeleteBuilder {
	return qb.Delete("participants").
		Where(squirrel.Eq{"id": id})
}

func getParticipantByUserAndEventQuery(userID, eventID string) squirrel.SelectBuilder {
	return qb.Select("id", "event_id", "user_id", "role", "gift_sent", "gift_sent_at", "created_at", "updated_at").
		From("participants").
		Where(squirrel.Eq{"user_id": userID, "event_id": eventID}).
		Limit(1)
}
