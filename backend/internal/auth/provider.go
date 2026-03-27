package auth

import (
	"context"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

type Provider interface {
	Config() *oauth2.Config
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}
