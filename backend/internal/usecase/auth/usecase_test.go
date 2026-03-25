package auth

import (
	"context"
	"errors"
	"testing"

	internalauth "secret-santa-backend/internal/auth"
	"secret-santa-backend/internal/entity"
)

type authRepoMock struct {
	getByOAuthIDFn func(ctx context.Context, oauthID string) (*entity.User, error)
	createFn       func(ctx context.Context, user entity.User) error
}

func (m authRepoMock) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return nil, nil
}
func (m authRepoMock) GetByOAuthID(ctx context.Context, oauthID string) (*entity.User, error) {
	return m.getByOAuthIDFn(ctx, oauthID)
}
func (m authRepoMock) Create(ctx context.Context, user entity.User) error {
	return m.createFn(ctx, user)
}

func TestLoginWithOAuthValidatesProviderID(t *testing.T) {
	uc := New(authRepoMock{})
	_, err := uc.LoginWithOAuth(context.Background(), internalauth.UserInfo{})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoginWithOAuthReturnsExistingUser(t *testing.T) {
	uc := New(authRepoMock{getByOAuthIDFn: func(ctx context.Context, oauthID string) (*entity.User, error) {
		return &entity.User{ID: "user-1"}, nil
	}})

	userID, err := uc.LoginWithOAuth(context.Background(), internalauth.UserInfo{ID: "oauth-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != "user-1" {
		t.Fatalf("unexpected userID: %s", userID)
	}
}

func TestLoginWithOAuthCreatesNewUser(t *testing.T) {
	called := false
	uc := New(authRepoMock{
		getByOAuthIDFn: func(ctx context.Context, oauthID string) (*entity.User, error) {
			return nil, errors.New("not found")
		},
		createFn: func(ctx context.Context, user entity.User) error {
			called = true
			if user.OAuthID != "oauth-1" || user.Name != "John" || user.Email != "john@example.com" || user.ID == "" {
				t.Fatalf("unexpected user: %+v", user)
			}
			return nil
		},
	})

	userID, err := uc.LoginWithOAuth(context.Background(), internalauth.UserInfo{ID: "oauth-1", Name: "John", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID == "" {
		t.Fatal("expected user id")
	}
	if !called {
		t.Fatal("repo create not called")
	}
}
