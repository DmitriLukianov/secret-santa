package scheduler

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type EventRepository interface {
	GetDueForDraw(ctx context.Context) ([]entity.Event, error)
	GetPendingDraws(ctx context.Context) ([]entity.Event, error)
}

type DrawUseCase interface {
	AutoDraw(ctx context.Context, eventID uuid.UUID) error
}

type DrawScheduler struct {
	eventRepo EventRepository
	drawUC    DrawUseCase
	log       *slog.Logger

	mu     sync.Mutex
	timers map[uuid.UUID]*time.Timer
}

func New(eventRepo EventRepository, drawUC DrawUseCase, log *slog.Logger) *DrawScheduler {
	return &DrawScheduler{
		eventRepo: eventRepo,
		drawUC:    drawUC,
		log:       log,
		timers:    make(map[uuid.UUID]*time.Timer),
	}
}

// Schedule registers a precise timer for the given event.
// Any previously registered timer for the same event is cancelled first.
func (s *DrawScheduler) Schedule(eventID uuid.UUID, drawAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cancelLocked(eventID)

	delay := time.Until(drawAt)
	if delay <= 0 {
		go s.runOne(context.Background(), eventID)
		return
	}

	t := time.AfterFunc(delay, func() {
		s.mu.Lock()
		delete(s.timers, eventID)
		s.mu.Unlock()
		s.runOne(context.Background(), eventID)
	})
	s.timers[eventID] = t
	s.log.Info("scheduler: draw scheduled",
		slog.String("event_id", eventID.String()),
		slog.String("in", delay.Round(time.Second).String()),
	)
}

// Cancel removes any pending timer for the given event.
func (s *DrawScheduler) Cancel(eventID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.timers[eventID]; ok {
		s.cancelLocked(eventID)
		s.log.Info("scheduler: draw cancelled", slog.String("event_id", eventID.String()))
	}
}

func (s *DrawScheduler) cancelLocked(eventID uuid.UUID) {
	if t, ok := s.timers[eventID]; ok {
		t.Stop()
		delete(s.timers, eventID)
	}
}

func (s *DrawScheduler) Start(ctx context.Context) {
	go func() {
		s.log.Info("draw scheduler started")
		s.loadPending(ctx)

		// Safety-net: catch draws missed while the server was down
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.log.Info("draw scheduler stopped")
				return
			case <-ticker.C:
				s.runDue(ctx)
			}
		}
	}()
}

func (s *DrawScheduler) loadPending(ctx context.Context) {
	events, err := s.eventRepo.GetPendingDraws(ctx)
	if err != nil {
		s.log.Error("scheduler: failed to load pending draws", slog.String("error", err.Error()))
		return
	}
	for _, e := range events {
		if e.DrawDate != nil {
			s.Schedule(e.ID, *e.DrawDate)
		}
	}
	s.log.Info("scheduler: pending draws loaded", slog.Int("count", len(events)))
}

func (s *DrawScheduler) runDue(ctx context.Context) {
	events, err := s.eventRepo.GetDueForDraw(ctx)
	if err != nil {
		s.log.Error("scheduler: failed to fetch due events", slog.String("error", err.Error()))
		return
	}
	for _, e := range events {
		go s.runOne(ctx, e.ID)
	}
}

func (s *DrawScheduler) runOne(ctx context.Context, eventID uuid.UUID) {
	if err := s.drawUC.AutoDraw(ctx, eventID); err != nil {
		s.log.Error("scheduler: auto draw failed",
			slog.String("event_id", eventID.String()),
			slog.String("error", err.Error()),
		)
	} else {
		s.log.Info("scheduler: auto draw completed",
			slog.String("event_id", eventID.String()),
		)
	}
}
