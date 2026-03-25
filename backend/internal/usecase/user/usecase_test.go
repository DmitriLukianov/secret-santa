package user

import (
	"context"
	"errors"
	"testing"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
)

type userRepoMock struct {
	createFn  func(ctx context.Context, user entity.User) error
	getByIDFn func(ctx context.Context, id string) (*entity.User, error)
	getAllFn  func(ctx context.Context) ([]entity.User, error)
	updateFn  func(ctx context.Context, id string, name, email *string) error
	deleteFn  func(ctx context.Context, id string) error
}

func (m userRepoMock) Create(ctx context.Context, user entity.User) error {
	return m.createFn(ctx, user)
}
func (m userRepoMock) GetByID(ctx context.Context, id string) (*entity.User, error) {
	return m.getByIDFn(ctx, id)
}
func (m userRepoMock) GetAll(ctx context.Context) ([]entity.User, error) {
	return m.getAllFn(ctx)
}
func (m userRepoMock) Update(ctx context.Context, id string, name, email *string) error {
	return m.updateFn(ctx, id, name, email)
}
func (m userRepoMock) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func TestUseCaseCreateValidatesInput(t *testing.T) {
	uc := New(userRepoMock{})

	if err := uc.Create(context.Background(), dto.CreateUserInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUseCaseCreateCallsRepo(t *testing.T) {
	called := false
	uc := New(userRepoMock{createFn: func(ctx context.Context, user entity.User) error {
		called = true
		if user.Name != "John" || user.Email != "john@example.com" {
			t.Fatalf("unexpected user: %+v", user)
		}
		return nil
	}})

	err := uc.Create(context.Background(), dto.CreateUserInput{Name: "John", Email: "john@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("repo not called")
	}
}

func TestUseCaseGetValidatesID(t *testing.T) {
	uc := New(userRepoMock{})
	if _, err := uc.Get(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUseCaseUpdateValidatesInput(t *testing.T) {
	uc := New(userRepoMock{})
	if err := uc.Update(context.Background(), "id", dto.UpdateUserInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUseCaseDeleteValidatesID(t *testing.T) {
	uc := New(userRepoMock{})
	if err := uc.Delete(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUseCaseGetAllReturnsRepoResult(t *testing.T) {
	uc := New(userRepoMock{getAllFn: func(ctx context.Context) ([]entity.User, error) {
		return []entity.User{{ID: "1"}}, nil
	}})

	users, err := uc.GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 || users[0].ID != "1" {
		t.Fatalf("unexpected users: %+v", users)
	}
}

func TestUseCaseCreatePropagatesRepoError(t *testing.T) {
	wantErr := errors.New("boom")
	uc := New(userRepoMock{createFn: func(ctx context.Context, user entity.User) error { return wantErr }})
	if err := uc.Create(context.Background(), dto.CreateUserInput{Name: "John", Email: "john@example.com"}); !errors.Is(err, wantErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
