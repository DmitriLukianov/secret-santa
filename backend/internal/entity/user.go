package entity

import (
	"time"

	"github.com/google/uuid"
)

// User — пользователь системы (сотрудник или клиент)
type User struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	OAuthID       string    `db:"oauth_id"`
	OAuthProvider string    `db:"oauth_provider"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// NewUser — фабрика пользователя (используется в usecase/auth)
func NewUser(name, email, oauthID, oauthProvider string) User {
	now := time.Now()
	return User{
		ID:            uuid.New(),
		Name:          name,
		Email:         email,
		OAuthID:       oauthID,
		OAuthProvider: oauthProvider,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
