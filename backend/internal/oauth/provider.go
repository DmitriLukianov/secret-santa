package oauth

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	"secret-santa-backend/internal/config" // ← важно
)

// UserInfo — общая структура для всех провайдеров
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

// Provider — главный интерфейс (точно по boilerplate)
type Provider interface {
	Config() *oauth2.Config
	GetAuthURL(state string) string
	GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error)
}

// New — фабрика провайдеров
func New(cfg *config.Config) (Provider, error) {
	switch cfg.OAuthProvider {
	case "github":
		return NewGitHubProvider(
			cfg.GithubClientID,
			cfg.GithubClientSecret,
			cfg.GithubRedirectURL,
		), nil

	// Здесь в будущем добавляем СДЭК ID, VK, Google и т.д.
	// case "sdek":
	//     return NewSdekProvider(cfg), nil
	// case "vk":
	//     return NewVKProvider(cfg), nil
	// case "google":
	//     return NewGoogleProvider(cfg), nil

	default:
		return nil, fmt.Errorf("unsupported oauth provider: %s", cfg.OAuthProvider)
	}
}
