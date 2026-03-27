package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// GetByID — запрос для получения пользователя по ID
func GetByID(id uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id})
}

// GetByOAuthID — запрос для поиска по OAuth
func GetByOAuthID(oauthID, oauthProvider string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{
			"oauth_id":       oauthID,
			"oauth_provider": oauthProvider,
		})
}

// GetByEmail — запрос по email
func GetByEmail(email string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email})
}

// GetAll — запрос для получения всех пользователей
func GetAll() squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users")
}

// Create — запрос на создание пользователя
func Create() squirrel.InsertBuilder {
	return qb.Insert("users").
		Columns("name", "email", "oauth_id", "oauth_provider").
		Suffix("RETURNING id, name, email, oauth_id, oauth_provider, created_at, updated_at")
}

// Update — запрос на частичное обновление
func Update(id uuid.UUID) squirrel.UpdateBuilder {
	return qb.Update("users").
		Where(squirrel.Eq{"id": id})
}

// Delete — запрос на удаление
func Delete(id uuid.UUID) squirrel.DeleteBuilder {
	return qb.Delete("users").
		Where(squirrel.Eq{"id": id})
}
