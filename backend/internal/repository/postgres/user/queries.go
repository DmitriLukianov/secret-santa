package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

var qb = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// GetByID возвращает запрос для получения пользователя по ID
func GetByID(id uuid.UUID) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id})
}

// GetByOAuthID возвращает запрос для поиска по OAuth-провайдеру
func GetByOAuthID(oauthID, oauthProvider string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{
			"oauth_id":       oauthID,
			"oauth_provider": oauthProvider,
		})
}

// GetByEmail возвращает запрос по email
func GetByEmail(email string) squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email})
}

// GetAll возвращает запрос на всех пользователей
func GetAll() squirrel.SelectBuilder {
	return qb.Select("id", "name", "email", "oauth_id", "oauth_provider", "created_at", "updated_at").
		From("users")
}

// Create возвращает INSERT с RETURNING — обязательно для получения данных из БД
func Create() squirrel.InsertBuilder {
	return qb.Insert("users").
		Columns("name", "email", "oauth_id", "oauth_provider").
		Suffix("RETURNING id, name, email, oauth_id, oauth_provider, created_at, updated_at")
}

// Update возвращает UPDATE с автоматическим обновлением updated_at
func Update(id uuid.UUID) squirrel.UpdateBuilder {
	return qb.Update("users").
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id})
}

// Delete возвращает DELETE builder
func Delete(id uuid.UUID) squirrel.DeleteBuilder {
	return qb.Delete("users").
		Where(squirrel.Eq{"id": id})
}
