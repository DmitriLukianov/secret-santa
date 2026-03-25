package participant

import (
	"context"
	"errors"
	"testing"

	"secret-santa-backend/internal/dto"
	"secret-santa-backend/internal/entity"
)

type participantRepoMock struct {
	getByEventFn func(ctx context.Context, eventID string) ([]entity.Participant, error)
	addFn        func(ctx context.Context, p entity.Participant) error
	deleteFn     func(ctx context.Context, id string) error
}

func (m participantRepoMock) GetByEvent(ctx context.Context, eventID string) ([]entity.Participant, error) {
	return m.getByEventFn(ctx, eventID)
}
func (m participantRepoMock) Add(ctx context.Context, p entity.Participant) error {
	return m.addFn(ctx, p)
}
func (m participantRepoMock) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

func TestParticipantAddValidatesInput(t *testing.T) {
	uc := New(participantRepoMock{})
	if err := uc.Add(context.Background(), dto.AddParticipantInput{}); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParticipantAddCallsRepo(t *testing.T) {
	called := false
	uc := New(participantRepoMock{addFn: func(ctx context.Context, p entity.Participant) error {
		called = true
		if p.EventID != "event-1" || p.UserID != "user-1" || p.ID == "" {
			t.Fatalf("unexpected participant: %+v", p)
		}
		return nil
	}})

	err := uc.Add(context.Background(), dto.AddParticipantInput{EventID: "event-1", UserID: "user-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("repo not called")
	}
}

func TestParticipantGetByEventValidatesEventID(t *testing.T) {
	uc := New(participantRepoMock{})
	if _, err := uc.GetByEvent(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParticipantDeleteValidatesID(t *testing.T) {
	uc := New(participantRepoMock{})
	if err := uc.Delete(context.Background(), ""); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestParticipantGetByEventPropagatesRepoError(t *testing.T) {
	wantErr := errors.New("boom")
	uc := New(participantRepoMock{getByEventFn: func(ctx context.Context, eventID string) ([]entity.Participant, error) {
		return nil, wantErr
	}})
	if _, err := uc.GetByEvent(context.Background(), "event-1"); !errors.Is(err, wantErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
