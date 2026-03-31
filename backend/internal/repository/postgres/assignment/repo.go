package assignment

import (
	"context"
	"fmt"

	"secret-santa-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a entity.Assignment) error {
	query := createAssignmentQuery().
		Values(a.ID, a.EventID, a.GiverID, a.ReceiverID, a.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) GetByEvent(ctx context.Context, eventID uuid.UUID) ([]entity.Assignment, error) {
	query := getAssignmentsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAssignments(rows)
}

func (r *Repository) DeleteByEvent(ctx context.Context, eventID uuid.UUID) error {
	query := deleteAssignmentsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

// FIXED: новая атомарная операция — вся жеребьёвка в одной транзакции
func (r *Repository) TransactionalDraw(ctx context.Context, eventID uuid.UUID, assignments []entity.Assignment, newStatus entity.EventStatus) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // в случае ошибки будет откат

	// 1. Удаляем старые назначения
	if err := r.deleteByEventTx(ctx, tx, eventID); err != nil {
		return fmt.Errorf("failed to delete old assignments: %w", err)
	}

	// 2. Создаём новые назначения
	for _, a := range assignments {
		if err := r.createTx(ctx, tx, a); err != nil {
			return fmt.Errorf("failed to create assignment: %w", err)
		}
	}

	// 3. Обновляем статус события
	if err := r.updateEventStatusTx(ctx, tx, eventID, newStatus); err != nil {
		return fmt.Errorf("failed to update event status: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Вспомогательные методы внутри транзакции
func (r *Repository) deleteByEventTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID) error {
	query := deleteAssignmentsByEventQuery(eventID.String())
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) createTx(ctx context.Context, tx pgx.Tx, a entity.Assignment) error {
	query := createAssignmentQuery().
		Values(a.ID, a.EventID, a.GiverID, a.ReceiverID, a.CreatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) updateEventStatusTx(ctx context.Context, tx pgx.Tx, eventID uuid.UUID, status entity.EventStatus) error {
	// Используем тот же query-builder, что и в event repo (можно вынести в общий helper при желании)
	// Здесь для минимальности используем прямой SQL
	_, err := tx.Exec(ctx, `
		UPDATE events 
		SET status = $1, updated_at = now() 
		WHERE id = $2`,
		status, eventID)
	return err
}
