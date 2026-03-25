package auth

import (
	"context"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	ID    string
	Email string
	Name  string
}

type Provider interface {
	Config() *oauth2.Config
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}
