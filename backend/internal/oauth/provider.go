package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	"secret-santa-backend/internal/config"
)

type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

type Provider interface {
	Config() *oauth2.Config
	GetAuthURL(state string) string
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}

func New(cfg *config.Config) (Provider, error) {
	switch cfg.OAuthProvider {
	case "github":
		return NewGitHubProvider(
			cfg.GithubClientID,
			cfg.GithubClientSecret,
			cfg.GithubRedirectURL,
		), nil

	default:
		return nil, fmt.Errorf("unsupported oauth provider: %s", cfg.OAuthProvider)
	}
}
