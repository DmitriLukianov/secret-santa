package notification

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createNotificationQuery() squirrel.InsertBuilder {
	return qb.Insert("notifications").
		Columns("user_id", "type", "payload")
}

func getNotificationsByUserQuery(userID uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "user_id", "type", "payload", "is_read", "created_at").
		From("notifications").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("created_at DESC")
}

func markAsReadQuery(id uuid.UUID) squirrel.UpdateBuilder {
	return qb.Update("notifications").
		Set("is_read", true).
		Where(squirrel.Eq{"id": id})
}

func markAllAsReadQuery(userID uuid.UUID) squirrel.UpdateBuilder {
	return qb.Update("notifications").
		Set("is_read", true).
		Where(squirrel.Eq{"user_id": userID, "is_read": false})
}
