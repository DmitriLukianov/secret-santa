package invitation

import "github.com/Masterminds/squirrel"

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func createInvitationQuery() squirrel.InsertBuilder {
	return qb.Insert("invitations").
		Columns("id", "event_id", "token", "expires_at", "created_by", "created_at", "updated_at")
}

func getInvitationByTokenQuery() squirrel.SelectBuilder {
	return qb.Select(
		"id", "event_id", "token", "expires_at",
		"created_by", "created_at", "updated_at",
	).
		From("invitations")
}
