package auth

import (
	"context"
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GitHubProvider struct {
	config *oauth2.Config
}

func NewGitHubProvider(clientID, clientSecret, redirectURL string) *GitHubProvider {
	return &GitHubProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     github.Endpoint,
			Scopes:       []string{"user:email"},
		},
	}
}

func (p *GitHubProvider) Config() *oauth2.Config {
	return p.config
}

func (p *GitHubProvider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func (p *GitHubProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (UserInfo, error) {
	client := p.config.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	var ghUser struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Login string `json:"login"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghUser); err != nil {
		return UserInfo{}, err
	}

	name := ghUser.Name
	if name == "" {
		name = ghUser.Login
	}

	email := ghUser.Email
	if email == "" {
		email = strconv.Itoa(ghUser.ID) + "@github.local"
	}

	return UserInfo{
		ID:       strconv.Itoa(ghUser.ID),
		Name:     name,
		Email:    email,
		Provider: "github", // ← теперь всегда передаём провайдер
	}, nil
}
