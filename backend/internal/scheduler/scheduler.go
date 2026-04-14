package scheduler

import (
	"context"
	"log/slog"
	"time"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
)

type EventRepository interface {
	GetDueForDraw(ctx context.Context) ([]entity.Event, error)
}

type DrawUseCase interface {
	AutoDraw(ctx context.Context, eventID uuid.UUID) error
}

type DrawScheduler struct {
	eventRepo EventRepository
	drawUC    DrawUseCase
	log       *slog.Logger
	interval  time.Duration
}

func New(eventRepo EventRepository, drawUC DrawUseCase, log *slog.Logger) *DrawScheduler {
	return &DrawScheduler{
		eventRepo: eventRepo,
		drawUC:    drawUC,
		log:       log,
		interval:  time.Minute,
	}
}

func (s *DrawScheduler) Start(ctx context.Context) {
	go func() {
		s.log.Info("draw scheduler started", slog.String("interval", s.interval.String()))

		// Запустить сразу при старте, не ждать первого тика
		s.runDraw(ctx)

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				s.log.Info("draw scheduler stopped")
				return
			case <-ticker.C:
				s.runDraw(ctx)
			}
		}
	}()
}

func (s *DrawScheduler) runDraw(ctx context.Context) {
	events, err := s.eventRepo.GetDueForDraw(ctx)
	if err != nil {
		s.log.Error("scheduler: failed to fetch events due for draw", slog.String("error", err.Error()))
		return
	}

	for _, event := range events {
		if err := s.drawUC.AutoDraw(ctx, event.ID); err != nil {
			s.log.Error("scheduler: auto draw failed",
				slog.String("event_id", event.ID.String()),
				slog.String("title", event.Title),
				slog.String("error", err.Error()),
			)
		} else {
			s.log.Info("scheduler: auto draw completed",
				slog.String("event_id", event.ID.String()),
				slog.String("title", event.Title),
			)
		}
	}
}
