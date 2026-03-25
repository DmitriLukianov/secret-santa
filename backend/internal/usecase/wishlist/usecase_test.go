package wishlist

import (
	"context"
	"errors"
	"testing"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
)

type wishlistRepoMock struct {
	createFn    func(ctx context.Context, w entity.Wishlist) error
	getByIDFn   func(ctx context.Context, id string) (*entity.Wishlist, error)
	getByUserFn func(ctx context.Context, userID string) ([]entity.Wishlist, error)
	updateFn    func(ctx context.Context, id string, title, description, link, imageURL, visibility *string) error
	deleteFn    func(ctx context.Context, id string) error
}

func (m wishlistRepoMock) Create(ctx context.Context, w entity.Wishlist) error {
	return m.createFn(ctx, w)
}
func (m wishlistRepoMock) GetByID(ctx context.Context, id string) (*entity.Wishlist, error) {
	return m.getByIDFn(ctx, id)
}
func (m wishlistRepoMock) GetByUser(ctx context.Context, userID string) ([]entity.Wishlist, error) {
	return m.getByUserFn(ctx, userID)
}
func (m wishlistRepoMock) Update(ctx context.Context, id string, title, description, link, imageURL, visibility *string) error {
	return m.updateFn(ctx, id, title, description, link, imageURL, visibility)
}
func (m wishlistRepoMock) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

func TestWishlistCreateValidatesUserID(t *testing.T) {
	uc := New(wishlistRepoMock{})
	if err := uc.Create(context.Background(), dto.CreateWishlistInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWishlistCreateCallsRepo(t *testing.T) {
	called := false
	uc := New(wishlistRepoMock{createFn: func(ctx context.Context, w entity.Wishlist) error {
		called = true
		if w.UserID != "user-1" || w.Title != "Wish" {
			t.Fatalf("unexpected wishlist: %+v", w)
		}
		return nil
	}})

	err := uc.Create(context.Background(), dto.CreateWishlistInput{UserID: "user-1", Title: "Wish"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("repo not called")
	}
}

func TestWishlistGetValidatesID(t *testing.T) {
	uc := New(wishlistRepoMock{})
	if _, err := uc.Get(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWishlistGetByUserValidatesUserID(t *testing.T) {
	uc := New(wishlistRepoMock{})
	if _, err := uc.GetByUser(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWishlistDeleteValidatesID(t *testing.T) {
	uc := New(wishlistRepoMock{})
	if err := uc.Delete(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestWishlistCreatePropagatesRepoError(t *testing.T) {
	wantErr := errors.New("boom")
	uc := New(wishlistRepoMock{createFn: func(ctx context.Context, w entity.Wishlist) error { return wantErr }})
	if err := uc.Create(context.Background(), dto.CreateWishlistInput{UserID: "user-1"}); !errors.Is(err, wantErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
