package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	OAuthID       string    `db:"oauth_id"`
	OAuthProvider string    `db:"oauth_provider"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func NewUser(name, email, oauthID, oauthProvider string) User {
	return User{
		Name:          name,
		Email:         email,
		OAuthID:       oauthID,
		OAuthProvider: oauthProvider,
	}
}
