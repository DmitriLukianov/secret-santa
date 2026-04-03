package invitation

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// createInvitationQuery — DB-first (убрали id, created_at, updated_at)
func createInvitationQuery() squirrel.InsertBuilder {
	return qb.Insert("invitations").
		Columns("event_id", "token", "expires_at", "created_by")
}

func getInvitationByTokenQuery() squirrel.SelectBuilder {
	return qb.Select(
		"id", "event_id", "token", "expires_at",
		"created_by", "created_at", "updated_at",
	).
		From("invitations")
}
