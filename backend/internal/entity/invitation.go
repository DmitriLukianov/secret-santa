package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID        uuid.UUID `db:"id"`
	EventID   uuid.UUID `db:"event_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedBy uuid.UUID `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// NewInvitation — чистый DB-first конструктор.
// ID, CreatedAt, UpdatedAt теперь генерирует PostgreSQL.
func NewInvitation(eventID, createdBy uuid.UUID, expiresIn time.Duration) Invitation {
	if expiresIn == 0 {
		expiresIn = 7 * 24 * time.Hour
	}

	return Invitation{
		EventID:   eventID,
		Token:     uuid.New().String(), // токен генерируем здесь (бизнес-логика)
		ExpiresAt: time.Now().Add(expiresIn),
		CreatedBy: createdBy,
		// ID, CreatedAt, UpdatedAt — будут заполнены БД
	}
}

func (i Invitation) IsValid() bool {
	return time.Now().Before(i.ExpiresAt)
}
