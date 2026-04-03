package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func GetByID(id uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id})
}

func GetByOAuthID(oauthID, oauthProvider string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{
			"oauth_id":       oauthID,
			"oauth_provider": oauthProvider,
		})
}

func GetByEmail(email string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email})
}

func GetAll() squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users")
}

func Create() squirrel.InsertBuilder {
	return qb.Insert("users").
		Columns("name", "email", "oauth_id", "oauth_provider").
		Suffix("RETURNING id, name, email, oauth_id, oauth_provider, created_at, updated_at")
}

func Update(id uuid.UUID) squirrel.UpdateBuilder {
	return qb.Update("users").
		Where(squirrel.Eq{"id": id})
}

func Delete(id uuid.UUID) squirrel.DeleteBuilder {
	return qb.Delete("users").
		Where(squirrel.Eq{"id": id})
}
