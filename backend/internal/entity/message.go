package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `db:"id"`
	EventID    uuid.UUID `db:"event_id"`
	SenderID   uuid.UUID `db:"sender_id"`
	ReceiverID uuid.UUID `db:"receiver_id"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
}

func NewMessage(eventID, senderID, receiverID uuid.UUID, content string) Message {
	return Message{
		EventID:    eventID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
	}
}
