package helpers

import (
	"fmt"
	"net/http"

	"secret-santa-backend/internal/definitions" // ← добавлен импорт

	"github.com/google/uuid"
)

// GetUserID — типизированный доступ к userID из контекста
func GetUserID(r *http.Request) (uuid.UUID, error) {
	val := r.Context().Value(definitions.UserIDKey)
	if val == nil {
		return uuid.Nil, fmt.Errorf("user id not found in context")
	}

	id, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user id type in context")
	}

	return id, nil
}
