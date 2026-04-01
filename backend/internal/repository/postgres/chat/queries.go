package chat

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createMessageQuery() squirrel.InsertBuilder {
	return qb.Insert("messages").
		Columns("id", "event_id", "sender_id", "receiver_id", "content", "created_at")
}

func getMessagesByPairQuery(eventID, user1ID, user2ID string) squirrel.SelectBuilder {
	return qb.Select("id", "event_id", "sender_id", "receiver_id", "content", "created_at").
		From("messages").
		Where(squirrel.Eq{"event_id": eventID}).
		Where(
			squirrel.Or{
				squirrel.And{
					squirrel.Eq{"sender_id": user1ID},
					squirrel.Eq{"receiver_id": user2ID},
				},
				squirrel.And{
					squirrel.Eq{"sender_id": user2ID},
					squirrel.Eq{"receiver_id": user1ID},
				},
			},
		).
		OrderBy("created_at ASC")
}
